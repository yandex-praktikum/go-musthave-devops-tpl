package server

import (
	"log"

	"github.com/i1l-ba/go-devops/internal/repository"
)

func initStore(config *Config) (syncChannel chan struct{}) {
	syncChannel = make(chan struct{}, 1)

	if config.StoreFilePath == "" {
		config.MetricsStore = repository.NewInMemoryStore()
	} else {
		fileStore, err := repository.NewFileStore(config.StoreFilePath, syncChannel)
		if err != nil {
			log.Fatalf("Failed to make file storage: %q", err)
		}

		config.MetricsStore = fileStore
	}

	return
}
