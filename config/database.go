package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetDatabaseURL constructs the database connection string based on environment variables
func GetDatabaseURL(cfg *Config) string {
	// Cloud SQL uses Unix domain socket when running on GCP
	host := os.Getenv("DB_HOST")

	// Check if we're connecting via Cloud SQL Proxy (Unix socket)
	if host[:9] == "/cloudsql" {
		return fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
			cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPassword)
	}

	return fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBPort)
}

// InitDB initializes the database connection
func InitDB(cfg *Config) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dsn := GetDatabaseURL(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set connection pool parameters
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set reasonable connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
