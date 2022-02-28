package workers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/itd27m01/go-metrics-service/internal/pkg/metrics"
	"github.com/itd27m01/go-metrics-service/internal/repository"
)

type ReporterConfig struct {
	ServerScheme   string `env:"SERVER_SCHEME" envDefault:"http"`
	ServerAddress  string `env:"ADDRESS"`
	ServerPath     string `env:"SERVER_PATH" envDefault:"/update/"`
	ServerTimeout  time.Duration
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

type ReportWorker struct {
	Cfg ReporterConfig
}

func (rw *ReportWorker) Run(ctx context.Context, mtr repository.Store) {
	reportTicker := time.NewTicker(rw.Cfg.ReportInterval)
	defer reportTicker.Stop()

	client := http.Client{
		Timeout: rw.Cfg.ServerTimeout,
	}

	serverURL := rw.Cfg.ServerScheme + "://" + rw.Cfg.ServerAddress + rw.Cfg.ServerPath

	for {
		select {
		case <-ctx.Done():
			return
		case <-reportTicker.C:
			SendReport(ctx, mtr, serverURL, &client)
			SendReportJSON(ctx, mtr, serverURL, &client)
			resetCounters(mtr)
		}
	}
}

func SendReport(ctx context.Context, mtr repository.Store, serverURL string, client *http.Client) {
	serverURL = strings.TrimSuffix(serverURL, "/")
	var stringifyMetricValue string
	for _, v := range mtr.GetMetrics() {
		if v.MType == metrics.GaugeMetricTypeName {
			stringifyMetricValue = fmt.Sprintf("%f", *v.Value)
		} else {
			stringifyMetricValue = fmt.Sprintf("%d", *v.Delta)
		}
		metricUpdateURL := fmt.Sprintf("%s/%s/%s/%s", serverURL, v.MType, v.ID, stringifyMetricValue)
		err := sendMetric(ctx, metricUpdateURL, client)
		if err != nil {
			log.Println(err)
		}
	}
}

func SendReportJSON(ctx context.Context, mtr repository.Store, serverURL string, client *http.Client) {
	serverURL = strings.TrimSuffix(serverURL, "/")
	updateURL := fmt.Sprintf("%s/", serverURL)
	for _, v := range mtr.GetMetrics() {
		err := sendMetricJSON(ctx, updateURL, client, v)
		if err != nil {
			log.Println(err)
		}
	}
}

func sendMetric(ctx context.Context, metricUpdateURL string, client *http.Client) error {
	log.Printf("Update metric: %s", metricUpdateURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, metricUpdateURL, nil)
	if err != nil {
		log.Println(err)

		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server response: %s", resp.Status)
	}
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func sendMetricJSON(ctx context.Context, serverURL string, client *http.Client, metric *metrics.Metric) error {
	log.Printf("Update metric: %s", metric.ID)

	body, err := metric.EncodeMetric()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server response: %s", resp.Status)
	}
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func resetCounters(mtr repository.Store) {
	mtr.ResetCounterMetric("PollCount")
}
