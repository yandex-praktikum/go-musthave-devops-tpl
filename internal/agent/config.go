package agent

//Не нашел сразу решение как сделать const/readonly структуру или мапу, по этому пока только так...
//Потом сделаю иначе
const (
	configPollInterval   = 2  //Seconds
	configReportInterval = 10 //Seconds
	configServerHost     = "127.0.0.1"
	configServerPort     = 8080
)
