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
	DBSSLMode          string
	DBConnectTimeout   string
	DBMaxConns         string
	DBMinConns         string
	DBMode             string // "local" or "production"
	UseIAMAuth         string // "true" or "false"
	AWSRegion          string
	OAIAPIKey          string
	GoogleClientSecret string
	WolframAppID       string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		DBMode:             getEnvOrDefault("DB_MODE", "production"),
		DBSSLMode:          getEnvOrDefault("DB_SSL_MODE", "enable"),
		DBConnectTimeout:   getEnvOrDefault("DB_CONNECT_TIMEOUT", "30"),
		DBMaxConns:         getEnvOrDefault("DB_MAX_CONNS", "10"),
		DBMinConns:         getEnvOrDefault("DB_MIN_CONNS", "2"),
		UseIAMAuth:         getEnvOrDefault("USE_IAM_AUTH", "true"),
		AWSRegion:          getEnvOrDefault("AWS_REGION", "us-east-1"),
		OAIAPIKey:          os.Getenv("OAI_API_KEY"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		WolframAppID:       os.Getenv("WOLFRAM_APP_ID"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
