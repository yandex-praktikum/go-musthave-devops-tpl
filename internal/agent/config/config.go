package config

import "time"

//Не нашел сразу решение как сделать const/readonly структуру или мапу, по этому пока только так...
//Потом сделаю иначе
const (
	ConfigClientRetryCount       = 3
	ConfigClientRetryWaitTime    = 10 * time.Second
	ConfigClientRetryMaxWaitTime = 90 * time.Second
	ConfigPollInterval           = 2  //Seconds
	ConfigReportInterval         = 10 //Seconds
	ConfigServerHost             = "127.0.0.1"
	ConfigServerPort             = 8080
)
