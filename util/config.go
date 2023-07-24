package util

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	MIDDLEWARE_TS_ADAPTER  = "ts_adapter"
	MIDDLEWARE_EVA_ADAPTER = "alfamart_eva"
)

var config Config

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	TS_URL        string `mapstructure:"TS_URL"`
	NEVA_URL      string `mapstructure:"NEVA_URL"`
	DB_User       string `mapstructure:"DB_USER"`
	DB_Pass       string `mapstructure:"DB_PASS"`
	DB_Address    string `mapstructure:"DB_ADDRESS"`
	DB_Port       string `mapstructure:"DB_PORT"`
	DB_Name       string `mapstructure:"DB_NAME"`
	GIN_MODE      string `mapstructure:"GIN_MODE"`
	LogFormat     string `mapstructure:"LOG_FORMAT"`
	LogOutput     string `mapstructure:"LOG_OUTPUT"`
	LogBeautify   bool   `mapstructure:"LOG_BEAUTIFY"`
	ENV           string `mapstructure:"ENV"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = viper.Unmarshal(&config)
	config.GIN_MODE = "release"
	return
}

func GetConfig() Config {
	return config
}
