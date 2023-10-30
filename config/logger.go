package config

type loggingConfig struct {
	EnableBark bool
	BarkConfig barkConfig
}

type barkConfig struct {
	ClientMode             string                 // Valid values are: SloggerOnly, DispatchToDb, DispatchToServer
	SloggerOnlyConfig      sloggerConfig          // Slogger mode config
	DispatchToDbConfig     dispatchToDbConfig     // Config for dispatching to the database
	DispatchToServerConfig dispatchToServerConfig // Config for dispatching to the bark server
}

type sloggerConfig struct {
	DefaultLogLevel string
}
type dispatchToDbConfig struct{}
type dispatchToServerConfig struct{}

func init() {
	config.Logging = loggingConfig{
		EnableBark: true,
		BarkConfig: barkConfig{
			ClientMode: "SloggerOnly",
			SloggerOnlyConfig: sloggerConfig{
				DefaultLogLevel: "INFO",
			},
			DispatchToDbConfig:     dispatchToDbConfig{},
			DispatchToServerConfig: dispatchToServerConfig{},
		},
	}
}

func initializeLoggingConfig() {
	config.Logging.EnableBark = envOrViperOrDefaultBool("logging.enableBark", config.Logging.EnableBark)
	config.Logging.BarkConfig.ClientMode = envOrViperOrDefaultString("logging.barkConfig.clientMode", config.Logging.BarkConfig.ClientMode)
}
