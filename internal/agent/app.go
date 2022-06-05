package agent

import (
	"log"
	"os"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"metrics/internal/agent/config"
	"metrics/internal/agent/requestHandler"
	"metrics/internal/agent/statsReader"
)

type AppHttp struct {
	isRun           bool
	startTime       time.Time
	lastRefreshTime time.Time
	lastUploadTime  time.Time
	client          *resty.Client
}

func (app *AppHttp) initHttpClient() {
	client := resty.New()

	client.
		SetRetryCount(config.ConfigClientRetryCount).
		SetRetryWaitTime(config.ConfigClientRetryWaitTime).
		SetRetryMaxWaitTime(config.ConfigClientRetryMaxWaitTime)

	app.client = client
}

func (app *AppHttp) Run() {
	var memStatistics statsReader.MemoryStatsDump
	signalChanel := make(chan os.Signal, 1)

	app.initHttpClient()
	app.startTime = time.Now()
	app.isRun = true

	tickerStatisticsRefresh := time.NewTicker(config.ConfigPollInterval * time.Second)
	tickerStatisticsUpload := time.NewTicker(config.ConfigReportInterval * time.Second)

	for app.isRun {
		select {
		case timeTickerRefresh := <-tickerStatisticsRefresh.C:
			log.Println("Refresh")
			app.lastRefreshTime = timeTickerRefresh
			memStatistics.Refresh()
		case timeTickerUpload := <-tickerStatisticsUpload.C:
			app.lastUploadTime = timeTickerUpload
			log.Println("Upload")

			err := requestHandler.MemoryStatsUpload(app.client, memStatistics)
			if err != nil {
				log.Println("Error!")
				log.Println(err)

				app.Stop()
			}
		case osSignal := <-signalChanel:
			switch osSignal {
			case syscall.SIGTERM:
				log.Println("syscall: SIGTERM")
			case syscall.SIGINT:
				log.Println("syscall: SIGINT")
			case syscall.SIGQUIT:
				log.Println("syscall: SIGQUIT")
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
