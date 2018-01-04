package sumo

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

//Settings ...
type Settings struct {
	DebugMode bool
	Port,
	SumoAddress,
	SumoAccessKey,
	SumoAccessId string
}

//ConfigKeys ...
var ConfigKeys *Settings

//InitConfig ...
func initConfig() *Settings {

	viper.SetDefault("ENVIRONMENT", "local")
	viper.SetDefault("APP_PORT", "8001")
	viper.SetDefault("DEBUG_Mode", false)

	if os.Getenv("ENVIRONMENT") == "local" {
		_, dirname, _, _ := runtime.Caller(0)
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(filepath.Dir(dirname))
		viper.ReadInConfig()

	} else {
		viper.AutomaticEnv()
	}

	return &Settings{
		Port:          viper.GetString("APP_PORT"),
		SumoAddress:   viper.GetString("SUMO_ADDRESS"),
		SumoAccessKey: viper.GetString("SUMO_ACCESS_KEY"),
		SumoAccessId:  viper.GetString("SUMO_ACCESS_ID"),
		DebugMode:     viper.GetBool("DEBUG_MODE"),
	}
}

//InitConfigSettings service configuration
func InitConfigSettings() {
	// SET config settings
	ConfigKeys = initConfig()
}
