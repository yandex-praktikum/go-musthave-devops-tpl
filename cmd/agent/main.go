package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"test_go/internal/agent"
	"time"
)

func readApplicationLogs() {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		agent.NewMonitor(2)
		//var id int
		//err := row.Scan(&id)

	}
}

func main() {

	//Start reader
	go readApplicationLogs()

	//Wait for exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//<-sigs
	log.Println("received", <-sigs)
	//log.Info("Received kill signal")
	//var tim
}
