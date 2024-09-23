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

	// Fetch the QtnPrice based on product and quantity
	if err := db.Where("city_product_id = ? AND quantity = ?", productID, qtn).First(&qtnPrice).Error; err != nil {
		return nil, fmt.Errorf("price not found for the specified quantity: %w", err)
	}

	var addresses []models.Address
	// Fetch only the available (non-assigned) addresses
	if err := db.Where("city_id = ? AND product_id = ? AND qtn_price_id = ? AND assigned = ?", cityID, productID, qtnPrice.ID, false).Find(&addresses).Error; err != nil {
		return nil, fmt.Errorf("address not found for the specified quantity: %w", err)
	}

	// If no available addresses are found, return an error
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no available addresses found for the given city, product, and quantity")
	}

	// Select a random address from the available ones
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random source
	randomIndex := r.Intn(len(addresses))                // Generate a random index
	randomAddress := addresses[randomIndex]

	// Mark the selected address as assigned and set the user who is assigned to this address
	randomAddress.Assigned = true
	randomAddress.AssignedUserID = &userID
	if err := db.Save(&randomAddress).Error; err != nil {
		return nil, fmt.Errorf("failed to assign address: %w", err)
	}

	// Prepare the address view to return
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
