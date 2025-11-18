package config

import (
	"os"

	helper "github.com/bncunha/erp-api/src/application/helpers"
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
	SMTP_HOST      string
	SMTP_PORT      int64
	SMTP_USERNAME  string
	SMTP_PASSWORD  string
	FRONTEND_URL   string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	smtpPort := helper.ParseInt64(os.Getenv("SMTP_PORT"))

	return &Config{
		DB_PASS:        os.Getenv("DB_PASS"),
		DB_USER:        os.Getenv("DB_USER"),
		DB_HOST:        os.Getenv("DB_HOST"),
		DB_PORT:        os.Getenv("DB_PORT"),
		DB_NAME:        os.Getenv("DB_NAME"),
		APP_ENV:        os.Getenv("APP_ENV"),
		NR_LICENSE_KEY: os.Getenv("NR_LICENSE_KEY"),
		SMTP_HOST:      os.Getenv("SMTP_HOST"),
		SMTP_PORT:      smtpPort,
		SMTP_USERNAME:  os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD:  os.Getenv("SMTP_PASSWORD"),
		FRONTEND_URL:   os.Getenv("FRONTEND_URL"),
	}, nil
}
