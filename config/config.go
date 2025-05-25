package config

import (
	"os"
)

type Config struct {
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	OAIAPIKey            string
	GoogleClientSecret   string
	StripeSecretKey      string
	StripePublishableKey string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:               os.Getenv("PGHOST"),
		DBPort:               os.Getenv("PGPORT"),
		DBUser:               os.Getenv("PGUSER"),
		DBPassword:           os.Getenv("PGPASSWORD"),
		DBName:               os.Getenv("PGDATABASE"),
		OAIAPIKey:            os.Getenv("OAI_API_KEY"),
		StripeSecretKey:      os.Getenv("STRIPE_SECRET_KEY"),
		StripePublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	}
}
