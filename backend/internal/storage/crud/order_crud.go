package crud

import (
	"Magaz/backend/internal/storage/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

// GetOrderByID fetches an order by its ID along with related entities such as User, City, Product, and potentially ReleasedBy and AddrToRelease.
func GetOrderByID(db *gorm.DB, orderID int) (*models.Order, error) { //TODO: Refactor to use crud template
	var order models.Order
	result := db.Preload("User").
		Preload("City").
		Preload("Product").
		Preload("ReleasedBy").
		Preload("AddrToRelease").
		First(&order, orderID) // 'First' adds a "WHERE id = ?" condition and limits the query to one row.

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Handle the case when no order is found
			log.Printf("No order found with ID: %d", orderID)
			return nil, result.Error
		}
		// Handle other potential errors
		log.Printf("Error retrieving order with ID: %d, error: %v", orderID, result.Error)
		return nil, result.Error
	}

	return &order, nil
}
