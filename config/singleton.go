package config

type masterConf struct {
	AppName  string
	Startup  startup
	Redis    redis
	Database database
	Logging  loggingConfig
}

var config masterConf

func Store() *masterConf {
	return &config
}

func InitializeConfig() {
	// Should the startup variables be printed or not
	determineLaunchInfoPrintSetting()

	initializeViper()
	initializeStartupConfig()
	initializeDatabaseConfig()
	initializeRedisConfig()
	initializeLoggingConfig()
}
