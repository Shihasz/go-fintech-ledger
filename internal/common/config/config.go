package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	// To add rest later
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // Look for config in this directory
	viper.SetConfigName("app") // Look for a file named "app.env"
	viper.SetConfigType("env") // The file type is "env"

	viper.AutomaticEnv() // Automatically override values with Environment Variables

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
