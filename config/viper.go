package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// vyper is the viper instance we will use
var vyper *viper.Viper

// NOTE: It's better not to use init function here

// initializeViper initializes the Viper variable and sets the config file name and search paths
func initializeViper() {
	vyper = viper.New()
	vyper.SetConfigName("config")
	vyper.SetConfigType("yaml")
	vyper.AddConfigPath("/.config/")
	vyper.AddConfigPath(".")

	err := readInViperConfig()
	if err != nil {
		fmt.Printf("E#1MQTGP - Viper could not be initialized. We can still read values from environment variables and defaults: %v", err)
	}
}

func readInViperConfig() error {
	if err := vyper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			errMessage := fmt.Sprintf("E#1MQTH0 - Config file not found: %v", err)
			fmt.Println(errMessage)
			return errors.New(errMessage)
		} else {
			// Config file was found but another error was produced
			errMessage := fmt.Sprintf("E#1MQTH2 - Config file found but another error occurred: %v", err)
			fmt.Println(errMessage)
			return errors.New(errMessage)
		}
	}
	return nil
}
