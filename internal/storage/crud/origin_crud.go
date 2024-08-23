package crud

import (
	"Magaz/internal/storage/models"
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

func GetAllCities(db *gorm.DB) ([]models.City, error) {
	var cities []models.City
	err := db.Find(&cities).Error
	return cities, err
}
