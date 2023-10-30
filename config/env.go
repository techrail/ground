package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultConfig     = "Default Value"
	configFileConfig  = "Config File"
	envVarConfigValue = "Environment Variable"
)

func envOrViperOrDefaultString(viperConfigKey string, defaultValue string) string {
	configKeyAsEnvVar := strings.ToUpper(strings.ReplaceAll(viperConfigKey, ".", "_"))
	value := os.Getenv(configKeyAsEnvVar)
	fetchedFrom := envVarConfigValue
	if value == "" {
		// Not in environment. Try the Viper value
		value = vyper.GetString(viperConfigKey)
		fetchedFrom = configFileConfig
		if value == "" {
			// Not found in viper config either
			value = defaultValue
			fetchedFrom = defaultConfig
		}
	}
	if printStartupConfigSetting {
		fmt.Printf("L#UJY8P - Value of `%v` [From %v] = `%v`\n", viperConfigKey, fetchedFrom, value)
	} else {
		fmt.Printf("L#18J4WN - Value of `%v` [From %v] = <Printing configuration is disabled>\n", viperConfigKey, fetchedFrom)
	}

	return value
}

func envOrViperOrDefaultInt64(viperConfigKey string, defaultValue int64) int64 {
	configKeyAsEnvVar := strings.ToUpper(strings.ReplaceAll(viperConfigKey, ".", "_"))
	value, err := strconv.ParseInt(os.Getenv(configKeyAsEnvVar), 0, 64)
	fetchedFrom := envVarConfigValue
	if err != nil {
		// Not in environment, or it was but not a valid integer value
		// Try the Viper value
		value = vyper.GetInt64(viperConfigKey)
		fetchedFrom = configFileConfig
		if value == int64(0) {
			// Not found in viper config either
			value = defaultValue
			fetchedFrom = defaultConfig
		}
	}

	if printStartupConfigSetting {
		fmt.Printf("L#UJZ1Q - Value of `%v` [From %v] = `%v`\n", viperConfigKey, fetchedFrom, value)
	} else {
		fmt.Printf("L#18J4Z3 - Value of `%v` [From %v] = <Printing configuration is disabled>\n", viperConfigKey, fetchedFrom)
	}
	return value
}

func envOrViperOrDefaultBool(viperConfigKey string, defaultValue bool) bool {
	configKeyAsEnvVar := strings.ToUpper(strings.ReplaceAll(viperConfigKey, ".", "_"))
	value, err := strconv.ParseBool(os.Getenv(configKeyAsEnvVar))
	fetchedFrom := envVarConfigValue
	if err != nil {
		// Not in environment, or it was but not a valid integer value
		// Try the Viper value
		val := vyper.Get(viperConfigKey)
		if val == nil {
			// Viper could not find the key in config (it is missing in the config file probably)
			value = defaultValue
			fetchedFrom = defaultConfig
		} else {
			value = vyper.GetBool(viperConfigKey)
			fetchedFrom = configFileConfig
		}
	}

	if printStartupConfigSetting {
		fmt.Printf("L#UJZUB - Value of `%v` [From %v] = `%v`\n", viperConfigKey, fetchedFrom, value)
	} else {
		fmt.Printf("L#18J4ZK - Value of `%v` [From %v] = <Printing configuration is disabled>\n", viperConfigKey, fetchedFrom)
	}
	return value
}

// envOrDefaultBoolean will check if an environment variable by the name of envVarName exists or not
// If the environment variable was set, then it returns the value as a boolean value or else returns the defaultValue
// NOTE: The parsing depends on strconv.ParseBool function. Look at that function's definition to figure out what
// qualifies as true and what qualifies as false. If the parsing fails, defaultValue is returned.
func envOrDefaultBoolean(envVarName string, defaultValue bool) bool {
	value := os.Getenv(envVarName)
	if value == "" {
		return defaultValue
	}

	retVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return retVal
}
