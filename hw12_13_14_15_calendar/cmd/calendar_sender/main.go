package main

import (
	"flag"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/logger"
	"os"

	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/service/sender"
)

var configFile string

func main() {
	flag.Parse()
	flag.StringVar(&configFile, "config", "configs/calendar.yaml", "Path to configuration file")

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}
	logg := logger.SetupLogger(cfg.Env)

	eventQueue, err := queue.NewQueue(cfg)
	if err != nil {
		logg.Error("failed to create queue: " + err.Error())
		os.Exit(1)
	}
	eventSender := sender.NewSender(*logg, eventQueue)
	eventSender.ReadMessages()
}
