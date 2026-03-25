package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	JWTSecret string
	DBUrl     string
	Debug     bool
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		AppPort:   getEnv("APP_PORT", "8000"),
		JWTSecret: getEnv("JWT_SECRET", ""),
		DBUrl:     getEnv("DATABASE_URL", ""),
		Debug:     getEnvAsBool("DEBUG", false),
	}

	validate(cfg)

	return cfg
}
