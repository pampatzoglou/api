package config

import (
	"os"
)

type ServerConfig struct {
	Port        string `envconfig:"SERVER_PORT" default:"8000"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"InfoLevel"`
	HealthPort  string `envconfig:"HEALTH_PORT" default:"9000"`
	MetricsPort string `envconfig:"METRICS_PORT" default:"9123"`
}

type DatabaseConfig struct {
	Connector string `envconfig:"MONGO_CONNECTOR" default:"mongodb://USER:PASS@localhost:27017/db"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port:     getEnv("PORT", "8000"),
			LogLevel: getEnv("LOG_LEVEL", "InfoLevel"),
		},
		Database: DatabaseConfig{
			Connector: getEnv("MONGO_CONNECTOR", "mongodb://USER:PASS@localhost:27017/db"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
