package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfig *envConfig

func InitializeEnvConfig() {
	EnvConfig = loadEnvVariables()
}

type envConfig struct {
	Host   string `mapstructure:"HOST"`
	Port   string `mapstructure:"PORT"`
	Env    string `mapstructure:"ENV"`
	PbHost string `mapstructure:"PB_HOST"`
	PbPort string `mapstructure:"PB_PORT"`
}

func loadEnvVariables() (config *envConfig) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file, %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}

	return
}
