package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	Port      string
	JWTSecret string
}

func Load() *Config {

	parentEnv := filepath.Join("..", ".env")
	err := godotenv.Load(parentEnv)
	if err != nil {
		log.Printf("No .env file found, relying on environment variables")
	}
	cfg := &Config{
		DBUrl:     getEnv("DATABASE_URL", ""),
		Port:      getEnv("APP_PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "default_secret"),
	}
	if cfg.DBUrl == "" {
		log.Fatal(("DATABASE_URL enviroment variable is required"))
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
