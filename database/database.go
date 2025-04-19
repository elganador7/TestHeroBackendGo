package database

import (
	"TestHeroBackendGo/config"
	"TestHeroBackendGo/models"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	// Configure SSL mode based on environment
	sslMode := cfg.DBSSLMode
	if cfg.DBMode == "production" {
		// Force SSL for production (Aurora)
		sslMode = "require"
	} else {
		// Disable SSL for local development
		sslMode = "disable"
	}

	// Get password (either from config or IAM token)
	password := cfg.DBPassword
	var err error
	if cfg.UseIAMAuth == "true" && cfg.DBMode == "production" {
		password, err = GenerateAuthToken(cfg)
		if err != nil {
			log.Fatal("Failed to generate IAM auth token:", err)
		}
	}

	// Construct the DSN with appropriate configuration
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, password, cfg.DBName, sslMode, cfg.DBConnectTimeout)

	// Parse connection pool settings
	maxConns, err := strconv.Atoi(cfg.DBMaxConns)
	if err != nil {
		log.Printf("Invalid DB_MAX_CONNS value, using default: %v", err)
		maxConns = 10
	}

	minConns, err := strconv.Atoi(cfg.DBMinConns)
	if err != nil {
		log.Printf("Invalid DB_MIN_CONNS value, using default: %v", err)
		minConns = 2
	}

	// Configure GORM with additional settings
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Adjust connection pool settings based on environment
	if cfg.DBMode == "production" {
		// More conservative settings for production
		sqlDB.SetMaxOpenConns(maxConns)
		sqlDB.SetMaxIdleConns(minConns)
		sqlDB.SetConnMaxLifetime(time.Hour)
	} else {
		// More relaxed settings for local development
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// Migrate the schema
	database.AutoMigrate(&models.Question{})

	DB = database
}
