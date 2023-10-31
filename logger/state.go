package logger

import (
	"github.com/techrail/bark/client"
	"sync"
)

const (
	Bark = "Bark"
)

type StateStruct struct {
	SelectedLogger string
	Initialized    bool
	Client         Logger
	sync.Mutex
}

var state *StateStruct

func init() {
	state = &StateStruct{
		SelectedLogger: Bark,
		Initialized:    false,
		Client:         client.NewSloggerClient(client.INFO),
	}
}

func State() *StateStruct {
	return state
}

/*
	Panic(string)
	Alert(string, bool)
	Error(string)
	Warn(string)
	Notice(string)
	Info(string)
	Debug(string)
*/

func Println(logmsg string) {
	state.Client.Println(logmsg)
}
