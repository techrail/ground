package logger

import "github.com/techrail/bark/client"

func NewLogger() Logger {
	BarkClient = client.NewSloggerClient(client.INFO)
	return BarkClient
}
