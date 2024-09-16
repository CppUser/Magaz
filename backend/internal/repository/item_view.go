package repository

import (
	"Magaz/backend/internal/storage/models"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type AddressView struct {
	ID          uint      `json:"id"`
	City        string    `json:"city"`
	Product     string    `json:"product"`
	Quantity    float32   `json:"quantity"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	AddedAt     time.Time `json:"added_at"`
	AddedBy     string    `json:"added_by"`
}

// TODO: Refactor it must return AddressView instead of model
func GetRandomAddress(db *gorm.DB, cityID uint, productID uint, qtn float32, userID int64) (*AddressView, error) {

	var qtnPrice models.QtnPrice

	if err := db.Where("city_product_id = ? AND quantity = ?", productID, qtn).First(&qtnPrice).Error; err != nil {
		return nil, fmt.Errorf("price not found for the specified quantity: %w", err)
	}

	var addresses []models.Address
	if err := db.Where("city_id = ? AND product_id = ? AND qtn_price_id = ?", cityID, productID, qtnPrice.ID).Find(&addresses).Error; err != nil {
		return nil, fmt.Errorf("address not found for the specified quantity: %w", err)
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses found for the given city, product, and quantity")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random source
	randomIndex := r.Intn(len(addresses))                // Generate a random index
	randomAddress := addresses[randomIndex]

	randomAddress.Assigned = true
	randomAddress.AssignedUserID = &userID
	if err := db.Save(&randomAddress).Error; err != nil {
		return nil, fmt.Errorf("failed to assign address: %w", err)
	}

	view := &AddressView{
		ID:          randomAddress.ID,
		City:        randomAddress.City.Name,
		Product:     randomAddress.Product.Name,
		Quantity:    qtn,
		Description: randomAddress.Description,
		Image:       randomAddress.Image,
		AddedAt:     randomAddress.AddedAt,
	}

	return view, nil

}
