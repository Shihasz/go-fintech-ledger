package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBSource      string `mapstructure:"DB_SOURCE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // Look for config in this directory
	viper.SetConfigName("app") // Look for a file named "app.env"
	viper.SetConfigType("env") // The file type is "env"

	viper.AutomaticEnv() // Automatically override values with Environment Variables

	_ = viper.BindEnv("SERVER_ADDRESS")
	_ = viper.BindEnv("DB_SOURCE")

	err = viper.ReadInConfig()
	if err != nil {
		// If file not found, rely on environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
