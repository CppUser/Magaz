package postgres

import (
	"Magaz/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	//Connect to the database
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	return nil, fmt.Errorf("failed to connect to the database: %w", err)
	//}
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		log.Printf("Failed to connect to the database, retrying in %v... (%d/%d)\n", 5*time.Second, i+1, 10)
		time.Sleep(5 * time.Second)
	}

	// Set the maximum number of open connections
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get the database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Println("Successfully connected to the database")
	return db, nil
}
