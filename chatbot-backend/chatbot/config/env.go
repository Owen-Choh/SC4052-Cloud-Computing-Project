package config

import (
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Port   string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Port: getEnv("BACKEND_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}