package config

import (
	"log"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Port                     string
	Default_Time             string
	Time_layout              string
	Timezone                 string
	JWTExpirationInSeconds   int64
	JWTSecret                string
	GEMINI_API_KEY           string
	API_FILE_EXPIRATION_HOUR int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Port:                     getEnv("BACKEND_PORT", "8080"),
		Default_Time:             "20 Mar 25 15:32 +0800",
		Time_layout:              "02 Jan 06 15:04 -0700",
		Timezone:                 "Asia/Singapore",
		JWTExpirationInSeconds:   getEnvInt("JWT_EXP", 3600*24*1),
		JWTSecret:                getEnv("JWT_SECRET", "should-have-jwt-secret-here"),
		GEMINI_API_KEY:           getEnv("GEMINI_API_KEY", ""),
		API_FILE_EXPIRATION_HOUR: getEnvInt("API_FILE_EXPIRATION_HOUR", 47),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Printf("Environment variable %s not set, using fallback value: %s", key, fallback)
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Environment variable %s not set, using fallback value: %d", key, fallback)
			return fallback
		}

		return intValue
	}

	log.Printf("Environment variable %s not set, using fallback value: %d", key, fallback)
	return fallback
}
