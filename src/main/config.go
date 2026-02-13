package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_PASS        string
	DB_USER        string
	DB_HOST        string
	DB_PORT        string
	DB_NAME        string
	APP_ENV        string
	NR_LICENSE_KEY string
	NR_ENABLED     bool
	BREVO_API_KEY  string
	FRONTEND_URL   string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		DB_PASS:        os.Getenv("DB_PASS"),
		DB_USER:        os.Getenv("DB_USER"),
		DB_HOST:        os.Getenv("DB_HOST"),
		DB_PORT:        os.Getenv("DB_PORT"),
		DB_NAME:        os.Getenv("DB_NAME"),
		APP_ENV:        os.Getenv("APP_ENV"),
		NR_LICENSE_KEY: os.Getenv("NR_LICENSE_KEY"),
		NR_ENABLED:     getBoolEnv("NR_ENABLED", true),
		BREVO_API_KEY:  os.Getenv("BREVO_API_KEY"),
		FRONTEND_URL:   os.Getenv("FRONTEND_URL"),
	}, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
