package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/i1l-ba/go-devops/cmd/server/cmd"
	"github.com/i1l-ba/go-devops/internal/server"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to parse command line arguments: %q", err)
	}

	metricsServerConfig := server.Config{
		ServerAddress: cmd.ServerAddress,
		StoreInterval: cmd.StoreInterval,
		Restore:       cmd.Restore,
		StoreFilePath: cmd.StoreFilePath,
	}
	if err := env.Parse(&metricsServerConfig); err != nil {
		log.Fatal(err)
	}

	metricsServer := server.MetricsServer{Cfg: &metricsServerConfig}

	metricsServer.Start(context.Background())
}
