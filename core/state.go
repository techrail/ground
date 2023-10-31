package core

import (
	"github.com/techrail/ground/logger"
)

// This file is supposed to contain the core state of the project which would allow the user to initialize the different components.

type groundState struct {
	LoggerState *logger.StateStruct
}

var state *groundState

func State() *groundState {
	SyncStates()
	return state
}

func SyncStates() {
	state.LoggerState = logger.State()
}

func InitBarkLogger() {
	logger.InitializeBarkLogger()
}
