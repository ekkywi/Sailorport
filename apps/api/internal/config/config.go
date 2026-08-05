package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port string
	AppEnv string
	Version string
}

func Load() Config {
	port := getenv("PORT", "8080")
	appEnv := getenv("APP_ENV", "development")
	version := getenv("APP_VERSION", "0.1.0")

	return Config{
		Port: port,
		AppEnv: appEnv,
		Version: version,
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func (c Config) PortInt() (int, error) {
	return strconv.Atoi(c.Port)
}