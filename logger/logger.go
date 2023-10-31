package logger

import "github.com/techrail/bark/client"

var BarkClient *client.Config

func InitializeBarkLogger() {
	BarkClient = client.NewSloggerClient(client.INFO)

	state.SelectedLogger = Bark
	state.Initialized = true
	state.Client = BarkClient
}
