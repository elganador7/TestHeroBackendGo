package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	GoogleClientSecret string
	OpenAIKey          string
	Port               string
	DatabaseURL        string
	Environment        string
	GCSBucketName      string
}

func LoadConfig() *Config {
	config := &Config{
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}

	config.DBHost = os.Getenv("DB_HOST")
	if config.DBHost == "" {
		log.Fatal("DB_HOST environment variable is required")
	}

	// Load OpenAI API key
	config.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	if config.OpenAIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Set port with default
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080" // Default port for Cloud Run
	}

	// Get environment
	config.Environment = os.Getenv("ENVIRONMENT")
	if config.Environment == "" {
		config.Environment = "development"
	}

	// Get GCS bucket name
	config.GCSBucketName = os.Getenv("GCS_BUCKET_NAME")
	if config.GCSBucketName == "" && config.Environment == "production" {
		log.Fatal("GCS_BUCKET_NAME environment variable is required in production")
	}

	return config
}
