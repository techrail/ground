package config

var printStartupConfigSetting = false

func determineLaunchInfoPrintSetting() {
	printStartupConfigSetting = envOrDefaultBoolean("PRINT_STARTUP_CONFIG", false)
}

func EnablePrintingConfig() {
	printStartupConfigSetting = true
}

func DisablePrintingConfig() {
	printStartupConfigSetting = false
}
