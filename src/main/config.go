package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_PASS string
	DB_USER string
	DB_HOST string
	DB_PORT string
	DB_NAME string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		DB_PASS: os.Getenv("DB_PASS"),
		DB_USER: os.Getenv("DB_USER"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_NAME: os.Getenv("DB_NAME"),
	}, nil
}