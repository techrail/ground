package channels

import (
	`github.com/techrail/bark/appRuntime`

	`github.com/techrail/ground/logger`
	`github.com/techrail/ground/typs/appError`
)

var ErrorStringChan chan string    // Channel to collect error strings from modules that can cause import cycles otherwise
var ErrorTypChan chan appError.Typ // Channel to collect appError types from modules that can cause import cycles otherwise

func init() {
	ErrorStringChan = make(chan string, 1000)
	ErrorTypChan = make(chan appError.Typ, 1000)

	go ProcessErrorStringChan()
	go ProcessErrorTypChan()
}

func ProcessErrorTypChan() {
	for {
		e := <-ErrorTypChan
		logger.Default(e.Error())
		if appRuntime.ShutdownRequested.Load() {
			return
		}
	}
}

func ProcessErrorStringChan() {
	for {
		e := <-ErrorStringChan
		logger.Default(e)
		if appRuntime.ShutdownRequested.Load() {
			return
		}
	}
}
