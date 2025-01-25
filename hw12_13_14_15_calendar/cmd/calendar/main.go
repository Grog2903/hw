package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/api/event"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage/sql"
	desc "github.com/Grog2903/hw/hw12_13_14_15_calendar/pkg/api/event/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalgrpc "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/server/http"
	isevent "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/service/event"
	memorystorage "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	logg := logger.SetupLogger(cfg.Env)

	var storage isevent.Storage
	switch cfg.Storage.Type {
	case "inMemory":
		storage = memorystorage.New()
	case "sql":
		sqlStorage := sqlstorage.New(nil)
		ctx := context.Background()
		if err := sqlStorage.Connect(ctx, *cfg); err != nil {
			logg.Error("failed connect to database: " + err.Error())
			os.Exit(1)
		}
		storage = sqlStorage
		defer sqlStorage.Close(ctx)
	}
	calendar := isevent.NewEventService(*logg, storage)
	controller := event.NewEventController(calendar)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServer.Port)) // :82
	if err != nil {
		slog.Error("failed to listen: %v", err)
	}

	grpcServer := internalgrpc.NewServer(*logg, *controller)
	err = grpcServer.Start(lis)
	if err != nil {
		slog.Error("grpc server error", err)
	}

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to dial server", err)
	}

	mux := runtime.NewServeMux()
	err = desc.RegisterCalendarHandler(context.Background(), mux, conn)
	if err != nil {
		slog.Error("failed to register calendar handler", err)
	}

	server := internalhttp.NewServer(*logg, *cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(mux); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
