package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DB_User string
	DB_Pass string
	DB_Host string
	DB_Port string
	DB_Name string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:    getEnv("PORT", "8080"),
		DB_User: getEnv("DB_USER", "app"),
		DB_Pass: getEnv("DB_PASSWORD", ""),
		DB_Host: getEnv("DB_HOST", "localhost"),
		DB_Port: getEnv("DB_PORT", "5432"),
		DB_Name: getEnv("DB_NAME", "db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
