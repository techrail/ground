package logger

import (
	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/valyala/fasthttp"
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

func LogWithContext(ctx *fasthttp.RequestCtx, msg string) {
	var opLog []string
	respondWithOplog := ctx.UserValue(customCtxKey.OpLogRequested)
	if respondWithOplog != nil {
		// Option was set in the context. Try to check its value
		valBool, ok := respondWithOplog.(bool)
		if ok && valBool {
			res := ctx.UserValue(customCtxKey.CtxOperationLogContent)
			if res == nil {
				errMsg := "E#1MVP0T - Value was nil. This was unexpected."
				Println(errMsg)
				opLog = []string{
					errMsg,
				}
			} else {
				// Try to assert
				if oprLog, typeAsserted := res.([]string); !typeAsserted {
					errMsg := "E#1MVP4B - Incorrect data format"
					Println(errMsg)
					opLog = []string{
						errMsg,
					}
				} else {
					opLog = oprLog
				}
			}

			opLog = append(opLog, msg)
			ctx.SetUserValue(customCtxKey.CtxOperationLogContent, opLog)
		}
	}

	// We might have to think another, better method to call here

	Println(msg)
}
