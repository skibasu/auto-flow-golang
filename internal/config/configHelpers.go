package config

import (
	"log"
	"os"
	"strconv"
)

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
