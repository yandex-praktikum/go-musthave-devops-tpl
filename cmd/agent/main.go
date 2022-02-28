package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/itd27m01/go-metrics-service/cmd/agent/cmd"
	"github.com/itd27m01/go-metrics-service/internal/workers"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to parse command line arguments: %q", err)
	}

	pollWorkerConfig := workers.PollerConfig{
		PollInterval: cmd.PollInterval,
	}
	if err := env.Parse(&pollWorkerConfig); err != nil {
		log.Fatal(err)
	}

	reportWorkerConfig := workers.ReporterConfig{
		ServerScheme:   "http",
		ServerAddress:  cmd.ServerAddress,
		ServerPath:     "/update/",
		ServerTimeout:  cmd.ServerTimeout,
		ReportInterval: cmd.ReportInterval,
	}
	if err := env.Parse(&reportWorkerConfig); err != nil {
		log.Fatal(err)
	}

	workers.Start(context.Background(), pollWorkerConfig, reportWorkerConfig)
}
