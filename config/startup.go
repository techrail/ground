package config

import (
	"fmt"
	"os"
	"strings"
)

type startup struct {
	Port            uint16
	EnableWebServer bool
	DebugMode       bool
	Env             string
}

func init() {
	// NOTE: Default config
	config.Startup = startup{
		Port:            8080,
		EnableWebServer: false,
		DebugMode:       false,
		Env:             "development",
	}
}

func initializeStartupConfig() {
	config.Startup.Port = uint16(envOrViperOrDefaultInt64("startup.port", int64(config.Startup.Port)))
	config.Startup.EnableWebServer = envOrViperOrDefaultBool("startup.enableWebServer", config.Startup.EnableWebServer)
	config.Startup.DebugMode = envOrViperOrDefaultBool("startup.debugMode", config.Startup.DebugMode)
	config.Startup.Env = envOrViperOrDefaultString("startup.env", config.Startup.Env)

	if !config.Startup.EnvIsProd() && !config.Startup.EnvIsUat() && !config.Startup.EnvIsLocalDev() {
		fmt.Printf("P#1DJ62U - NOT AN ACCEPTABLE ENVIRONMENT TO EXECUTE IN: %v\n", config.Startup.Env)
		os.Exit(1)
	}
}

func (s *startup) EnvIsProd() bool {
	if strings.ToUpper(s.Env) == "PROD" || strings.ToUpper(s.Env) == "PRODUCTION" {
		return true
	}
	return false
}

func (s *startup) EnvIsUat() bool {
	if strings.ToUpper(s.Env) == "UAT" {
		return true
	}
	return false
}

// EnvIsLocalDev checks if the environment is for local development
func (s *startup) EnvIsLocalDev() bool {
	if strings.ToUpper(s.Env) == "DEV" || strings.ToUpper(s.Env) == "DEVELOPMENT" {
		return true
	}
	return false
}
