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
	Secret  string
}

func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := Config{
		AppPort: getEnv("APP_PORT"),
		DBUrl:   getEnv("DATABASE_URL"),
		Debug:   getEnvAsBool("DEBUG"),
		Secret:  getEnv("JWT_SECRET"),
	}

	return cfg
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return "SeCrET"
}

func getEnvAsBool(key string) bool {
	valueStr := getEnv(key)
	if valueStr == "" {
		return false
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false
	}

	return value
}
