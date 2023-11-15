package config

type tym struct {
	Timezone string
}

func init() {
	// NOTE: Default config
	config.Time = tym{
		Timezone: "UTC",
	}
}

func initializeTimeConfig() {
	config.Time.Timezone = envOrViperOrDefaultString("time.timezone", config.Time.Timezone)
}
