package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SerwerAddress string
	JWTSecret     string
	DBUrl         string
	Debug         bool
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		SerwerAddress: getEnv("SERVER_ADDRESS", ":8000"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		DBUrl:         getEnv("DATABASE_URL", ""),
		Debug:         getEnvAsBool("DEBUG", false),
	}

	validate(cfg)

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return fallback
	}

	return value
}

func validate(cfg *Config) {
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	if cfg.DBUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

}
