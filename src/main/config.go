package config

import (
	"os"

	"github.com/joho/godotenv"
)

var COMPANY_TEST_ID = 1

type Config struct {
	DB_PASS string
	DB_HOST string
	DB_PORT string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB_PASS: os.Getenv("DB_PASS"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
	}, nil
}