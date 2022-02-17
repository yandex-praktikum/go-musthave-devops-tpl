package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test_go/internal/agent"
	"time"
)

func readApplicationLogs(cnt int) {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		agent.NewMonitor(2, addCounterChan, cnt)
		//var id int
		//err := row.Scan(&id)

	}
}

var counter int
var addCounterChan chan int
var readCounterChan chan int

//func AddCounter(ch chan int) {
//	ch <- 1
//}

func main() {
	addCounterChan = make(chan int, 100)
	readCounterChan = make(chan int, 100)

	counter = 0
	m := make(map[string]int)

	go func() {
		for {
			select {
			case val := <-addCounterChan:
				counter += val
				m["cnt"] = counter

				readCounterChan <- counter
				fmt.Printf("Count %d \n", counter)
			}
		}
	}()
	//Start reader
	go readApplicationLogs(m["cnt"])
	//Wait for exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//<-sigs
	log.Println("received", <-sigs)
	//log.Info("Received kill signal")
	//var tim
}
