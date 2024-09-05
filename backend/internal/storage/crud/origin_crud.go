package crud

import (
	"Magaz/backend/internal/storage/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func CreateCity(db *gorm.DB, city *models.City) error {
	return db.Create(city).Error
}

func GetCityByID(db *gorm.DB, id uint) (*models.City, error) {
	var city models.City
	err := db.Where("id = ?", id).First(&city).Error
	return &city, err
}

func GetCityIDByName(db *gorm.DB, name string) (uint, error) {
	var city models.City
	if err := db.Where("name = ?", name).First(&city).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("city not found: %w", err)
		}
		return 0, fmt.Errorf("failed to find city: %w", err)
	}

	return city.ID, nil
}

func GetAllCities(db *gorm.DB) ([]models.City, error) {
	var cities []models.City
	err := db.Find(&cities).Error
	return cities, err
}
