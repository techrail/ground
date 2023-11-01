package logger

import "github.com/techrail/bark/client"

var BarkClient *client.Config

func InitializeBarkLogger() {
	BarkClient = client.NewSloggerClient(client.INFO)
	state = &StateStruct{
		SelectedLogger: Bark,
		Initialized:    true,
		Client:         BarkClient,
	}
}
