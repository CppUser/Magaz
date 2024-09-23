package crud

import (
	"Magaz/backend/internal/storage/models"
	"log"

	"gorm.io/gorm"
)

func GetAvailableAddresses(db *gorm.DB, cityID, productID uint, quantity float32) ([]models.Address, error) {
	var qtnPrice models.QtnPrice

	// Fetch the correct QtnPriceID based on the quantity provided
	err := db.Where("city_product_id IN (SELECT id FROM city_products WHERE city_id = ? AND product_id = ?) AND quantity <= ?", cityID, productID, quantity).
		Order("quantity DESC"). // Prefer higher quantity match
		First(&qtnPrice).Error
	if err != nil {
		log.Printf("Failed to find matching QtnPrice for quantity: %v", err)
		return nil, err
	}

	//TODO: refactor use template
	//addresses, err := crud.GetAllWithAssociations[models.Address](
	//	db.Where("city_id = ? AND product_id = ? AND qtn_price_id = ? AND released = ? AND assigned = ?", cityID, productID, qtnPrice.ID, false, false),
	//	"City", "Product", "QtnPrice", // Associations to preload
	//)

	// Once the QtnPrice is found, fetch the available addresses
	var addresses []models.Address
	result := db.Where("city_id = ? AND product_id = ? AND qtn_price_id = ? AND released = ? AND assigned = ?", cityID, productID, qtnPrice.ID, false, false).
		Find(&addresses)

	if result.Error != nil {
		log.Printf("Failed to find available addresses: %v", result.Error)
		return nil, result.Error
	}

	return addresses, nil
}

func GetAddressByID(db *gorm.DB, addressID uint) (*models.Address, error) {
	var address models.Address
	if err := db.First(&address, addressID).Error; err != nil {
		log.Printf("Failed to find address: %v", err)
		return nil, err
	}
	return &address, nil
}
