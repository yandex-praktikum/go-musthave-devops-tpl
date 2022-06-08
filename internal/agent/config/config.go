package config

import "time"

const (
	ConfigClientRetryCount       = 1
	ConfigClientRetryWaitTime    = 10 * time.Second
	ConfigClientRetryMaxWaitTime = 90 * time.Second
	ConfigPollInterval           = 2  //Seconds
	ConfigReportInterval         = 10 //Seconds
	ConfigServerHost             = "127.0.0.1"
	ConfigServerPort             = 8080
)
