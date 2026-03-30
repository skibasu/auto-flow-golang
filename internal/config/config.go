package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBUrl   string
	Debug   bool
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		AppPort: GetEnv("APP_PORT"),
		DBUrl:   GetEnv("DATABASE_URL"),
		Debug:   getEnvAsBool("DEBUG"),
	}

	return cfg
}

func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

func getEnvAsBool(key string) bool {
	valueStr := GetEnv(key)
	if valueStr == "" {
		return false
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false
	}

	return value
}
