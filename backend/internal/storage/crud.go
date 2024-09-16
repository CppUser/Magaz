package crud

import (
	"errors"
	"gorm.io/gorm"
)

// Model interface that all models should implement
type Model interface {
	TableName() string
}

// Create a new record in the database
func Create[T any](db *gorm.DB, model *T) error {
	return db.Create(model).Error
}

// Get a record by ID from the database
func Get[T any, I comparable](db *gorm.DB, id I) (*T, error) {
	var model T
	if err := db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

// Get all records from the database
func GetAll[T any](db *gorm.DB) ([]T, error) {
	var models []T
	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

// Get a record by ID from the database with preloaded associations
func GetWithAssociations[T any](db *gorm.DB, id uint, associations ...string) (*T, error) {
	var model T
	query := db.First(&model, id)
	for _, assoc := range associations {
		query = query.Preload(assoc)
	}
	if err := query.Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func GetAllWithAssociations[T any](db *gorm.DB, associations ...string) ([]T, error) {
	var models []T
	query := db.Find(&models)
	for _, assoc := range associations {
		query = query.Preload(assoc)
	}
	if err := query.Error; err != nil {
		return nil, err
	}
	return models, nil
}

// Update a record in the database
func Update[T any](db *gorm.DB, model T) error {
	if db.Model(model).Updates(model).RowsAffected == 0 {
		return errors.New(" not found")
	}
	return nil
}

// Delete a record by ID from the database
func Delete[T any](db *gorm.DB, id uint) error {
	var model T
	if err := db.Delete(&model, id).Error; err != nil {
		return err
	}
	return nil
}
