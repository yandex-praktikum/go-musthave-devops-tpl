package requestHandler

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
	"metrics/internal/agent/config"
	"metrics/internal/agent/statsReader"
	"reflect"
	"strings"
)

func oneStatUpload(httpClient *resty.Client, statType string, statName string, statValue string) error {
	resp, err := httpClient.R().
		SetPathParams(map[string]string{
			"host":  config.ConfigServerHost,
			"port":  fmt.Sprintf("%v", config.ConfigServerPort),
			"type":  statType,
			"name":  statName,
			"value": statValue,
		}).
		SetHeader("Content-Type", "text/plain").
		Post("http://{host}:{port}/update/{type}/{name}/{value}")

	if err != nil {
		fmt.Println(err)
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New("HTTP Status != 200")
	}

	return nil
}

func MemoryStatsUpload(httpClient *resty.Client, memoryStats statsReader.MemoryStatsDump) error {
	reflectMemoryStats := reflect.ValueOf(memoryStats)
	typeOfMemoryStats := reflectMemoryStats.Type()
	errorGroup := new(errgroup.Group)

	for i := 0; i < reflectMemoryStats.NumField(); i++ {
		statName := typeOfMemoryStats.Field(i).Name
		statValue := fmt.Sprintf("%v", reflectMemoryStats.Field(i).Interface())
		statType := strings.Split(typeOfMemoryStats.Field(i).Type.String(), ".")[1]

		errorGroup.Go(func() error {
			return oneStatUpload(httpClient, statType, statName, statValue)
		})
	}

	err := errorGroup.Wait()
	return err
}
