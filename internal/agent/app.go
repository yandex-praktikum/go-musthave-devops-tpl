package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"metrics/internal/agent/requestHandler"
	"metrics/internal/agent/statsReader"
	"os"
	"syscall"
	"time"
)

type AppRunner interface {
	Run()
	Stop()
	IsRun() bool
}

type AppHttp struct {
	isRun           bool
	startTime       time.Time
	lastRefreshTime time.Time
	lastUploadTime  time.Time
	client          *resty.Client
}

func (app *AppHttp) Run() {
	var memStatistics statsReader.MemoryStatsDump
	signalChanel := make(chan os.Signal, 1)

	app.client = resty.New()
	app.startTime = time.Now()
	app.isRun = true

	tickerStatisticsRefresh := time.NewTicker(configPollInterval * time.Second)
	tickerStatisticsUpload := time.NewTicker(configReportInterval * time.Second)

	for {
		select {
		case timeTickerRefresh := <-tickerStatisticsRefresh.C:
			fmt.Println("Refresh")
			app.lastRefreshTime = timeTickerRefresh
			memStatistics.Refresh()
		case timeTickerUpload := <-tickerStatisticsUpload.C:
			app.lastUploadTime = timeTickerUpload
			fmt.Println("Upload")

			err := requestHandler.MemoryStatsUpload(memStatistics)
			if err != nil {
				app.Stop()
			}
		case osSignal := <-signalChanel:
			switch osSignal {
			case syscall.SIGTERM:
				fmt.Println("syscall: SIGTERM")
			case syscall.SIGINT:
				fmt.Println("syscall: SIGINT")
			case syscall.SIGQUIT:
				fmt.Println("syscall: SIGQUIT")
			}
			app.Stop()
		}
	}
}

func (app *AppHttp) Stop() {
	app.isRun = false
}

func (app *AppHttp) IsRun() bool {
	return app.isRun
}
