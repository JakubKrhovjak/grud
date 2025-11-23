package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string         `mapstructure:"env"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Grpc     GrpcConfig     `mapstructure:"grpc"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type GrpcConfig struct {
	Port string `mapstructure:"port"`
}

func Load() (*Config, error) {
	// Get environment from APP_ENV, default to "local"
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// Set up Viper
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
