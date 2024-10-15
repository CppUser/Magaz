package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
	"user/internal/config"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	var db *gorm.DB
	var err error

	// Retry connecting to the database up to 10 times
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Connection successful, break the retry loop
			break
		}

		// Log the error and sleep before retrying
		log.Printf("Failed to connect to the database: %v. Retrying in %v... (%d/%d)\n", err, 5*time.Second, i+1, 10)
		time.Sleep(5 * time.Second)
	}

	// If still no connection after retries, return the error
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database after retries: %w", err)
	}

	// Get the SQL database object and configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get the database connection: %w", err)
	}

	// Set the maximum number of open connections
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// Set the maximum number of idle connections
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// Set the maximum lifetime of a connection
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Println("Successfully connected to the database")
	return db, nil
}
