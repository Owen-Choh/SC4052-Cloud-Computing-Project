package config

import (
	"log"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	FrontendDomain           string
	Port                     string
	DATABASE_PATH            string
	FILES_PATH               string
	Default_Time             string
	Time_layout              string
	Timezone                 string
	JWTExpirationInSeconds   int64
	JWTSecret                string
	API_FILE_EXPIRATION_HOUR int64
	GEMINI_API_KEY           string
	MODEL_NAME               string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		FrontendDomain:           getEnv("FrontendDomain", "http://localhost:5173"),
		Port:                     getEnv("BACKEND_PORT", "8080"),
		DATABASE_PATH:            getEnv("DATABASE_PATH", "./database_files/chatbot.db"),
		FILES_PATH:               getEnv("FILES_PATH", "database_files/uploads/"),
		Default_Time:             getEnv("Default_Time", "20 Mar 25 15:32 +0800"),
		Time_layout:              getEnv("Time_layout", "02 Jan 06 15:04 -0700"),
		Timezone:                 getEnv("Timezone", "Asia/Singapore"),
		JWTExpirationInSeconds:   getEnvInt("JWT_EXP_SECONDS", 3600*24*1),
		JWTSecret:                getEnv("JWT_SECRET", "should-have-jwt-secret-here"),
		API_FILE_EXPIRATION_HOUR: getEnvInt("API_FILE_EXPIRATION_HOUR", 47),
		MODEL_NAME:               getEnv("MODEL_NAME", "gemini-2.0-flash-thinking-exp-01-21"),
		GEMINI_API_KEY:           getOSEnv("GEMINI_API_KEY", ""),
	}
}

func getOSEnv(key string, fallback string) string {
	value := os.Getenv("key") // Get API key from environment variable
	if value == "" {
		log.Printf("OS Environment variable %s not set, using fallback value: %s", key, fallback)
		return fallback
	}

	return value
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
