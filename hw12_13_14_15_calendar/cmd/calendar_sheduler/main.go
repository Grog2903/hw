package main

import (
	"flag"
	"os"

	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/service/scheduler"
	memorystorage "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage/sql"
	"golang.org/x/net/context"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yaml", "Path to configuration file")
	flag.Parse()

	code := run()
	os.Exit(code)
}

func run() int {
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}
	logg := logger.SetupLogger(cfg.Env)

	var storage scheduler.Storage

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

	eventQueue, err := queue.NewQueue(cfg)
	if err != nil {
		logg.Error("failed to create queue: " + err.Error())
		return 1
	}

	eventScheduler := scheduler.NewScheduler(*logg, storage, eventQueue)
	eventScheduler.Start(context.Background(), cfg.Scheduler.LaunchFrequency)

	return 0
}
