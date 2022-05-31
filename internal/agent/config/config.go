package config

//Не нашел сразу решение как сделать const/readonly структуру или мапу, по этому пока только так...
//Потом сделаю иначе
const (
	ConfigPollInterval   = 2 //Seconds
	ConfigReportInterval = 4 //Seconds
	ConfigServerHost     = "127.0.0.1"
	ConfigServerPort     = 8080
)
