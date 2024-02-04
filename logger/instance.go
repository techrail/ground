package logger

import "github.com/techrail/bark/client"

func NewSloggerClient() Logger {
	return client.NewSloggerClient(client.INFO)
}

func NewBarkClient(remoteServerUrl, defaultLogLvl, svcName, sessName string, enableSlog bool, enableBulkSend bool) Logger {
	return client.NewClient(remoteServerUrl, defaultLogLvl, svcName, sessName, enableSlog, enableBulkSend)
}

func NewEmbeddedServerBarkClient(dbUrl, defaultLogLvl, svcName, sessName string, enableSlog bool) Logger {
	return client.NewClientWithServer(dbUrl, defaultLogLvl, svcName, sessName, enableSlog)
}
