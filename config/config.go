package config

import (
	"os"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	OAIAPIKey          string
	GoogleClientSecret string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     os.Getenv("PRGHOST"),
		DBPort:     os.Getenv("PGPORT"),
		DBUser:     os.Getenv("PGUSER"),
		DBPassword: os.Getenv("PGPASSWORD"),
		DBName:     os.Getenv("PGDATABASE"),
		OAIAPIKey:  os.Getenv("OAI_API_KEY"),
	}
}
