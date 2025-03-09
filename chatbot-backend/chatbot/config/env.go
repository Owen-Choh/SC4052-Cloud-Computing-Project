package config

import (
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Port   string
	Default_Time string
	Timezone string
	JWTExpirationInSeconds int64
	JWTSecret string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Port: getEnv("BACKEND_PORT", "8080"),
		Default_Time: "2025-01-01 00:00:00",
		Timezone: "Asia/Singapore",
		JWTExpirationInSeconds: getEnvInt("JWT_EXP", 3600*24*7),
		JWTSecret: getEnv("JWT_SECRET", "should-have-jwt-secret-here"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return intValue
	}

	return fallback
}