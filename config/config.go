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
	Host           string   `mapstructure:"HOST"`
	Port           string   `mapstructure:"PORT"`
	Env            string   `mapstructure:"ENV"`
	PbHost         string   `mapstructure:"PB_HOST"`
	PbPort         string   `mapstructure:"PB_PORT"`
	AdminUser      string   `mapstructure:"ADMIN_USER"`
	AdminPass      string   `mapstructure:"ADMIN_PASS"`
	AllowedOrigins []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	// UploadMaxBytes configures the maximum allowed uploaded file size in bytes. If unset, defaults to 10MB.
	UploadMaxBytes int64 `mapstructure:"UPLOAD_MAX_BYTES"`
}

func loadEnvVariables() (config *envConfig) {
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "PROD")
	viper.SetDefault("PB_HOST", "localhost")
	viper.SetDefault("PB_PORT", "8090")
	viper.SetDefault("ADMIN_USER", "admin@hosting.test")
	viper.SetDefault("ADMIN_PASS", "pleaseChange123!")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", nil)
	// default upload max bytes: 10 MB
	viper.SetDefault("UPLOAD_MAX_BYTES", 10<<20)

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No .env file found, using default values and environment variables.")
		} else {
			log.Fatalf("Error reading .env file, %s", err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}

	return
}
