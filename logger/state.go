package logger

import (
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
	InitializeBarkLogger()
}

func EnableDebug() {
	if state.Initialized && state.SelectedLogger == Bark {
		BarkClient.EnableDebugLogs()
	}
}

func DisableDebug() {
	if state.Initialized && state.SelectedLogger == Bark {
		BarkClient.DisableDebugLogs()
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

func Panic(msg string) {
	state.Client.Println(msg)
}

func Alert(msg string) {
	state.Client.Alert(msg, false)
}

func AlertWait(msg string) {
	state.Client.Alert(msg, true)
}

func Error(msg string) {
	state.Client.Error(msg)
}

func Warn(msg string) {
	state.Client.Warn(msg)
}

func Notice(msg string) {
	state.Client.Notice(msg)
}

func Info(msg string) {
	state.Client.Info(msg)
}

func Debug(msg string) {
	state.Client.Debug(msg)
}

func Println(msg string) {
	state.Client.Println(msg)
}
