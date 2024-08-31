package repository

import (
	"Magaz/backend/internal/storage/models"
	"gorm.io/gorm"
)

// ProductItem represents a specific quantity and price of a product.
type ProductItem struct {
	Quantity  float32
	Price     float32
	Available float32
}

// ProductView is used for displaying products in the UI.
type ProductView struct {
	Name  string
	Total float32 // This represents the total available quantity
	Items []ProductItem
}

// CityWithProductsView represents a city and its associated products for display.
type CityWithProductsView struct {
	City     string
	Products []ProductView
}

func FetchCityProducts(db *gorm.DB) ([]CityWithProductsView, error) {
	var cityProducts []models.CityProduct

	// Preload associated City, Product, and ProductPrices
	err := db.Preload("City").Preload("Product").Preload("ProductPrices").Find(&cityProducts).Error
	if err != nil {
		return nil, err
	}

	// Group data by city
	cityMap := make(map[uint]*CityWithProductsView)

	for _, cp := range cityProducts {
		city := cp.City
		product := cp.Product

		// Initialize city in the map if not already present
		if _, exists := cityMap[city.ID]; !exists {
			cityMap[city.ID] = &CityWithProductsView{
				City:     city.Name,
				Products: []ProductView{},
			}
		}

		// Convert ProductPrices to ProductItems
		var productItems []ProductItem
		for _, pp := range cp.ProductPrices {
			productItems = append(productItems, ProductItem{
				Quantity:  pp.Quantity,
				Price:     pp.Price,
				Available: cp.TotalQuantity,
			})
		}

		// Add product under the city
		cityMap[city.ID].Products = append(cityMap[city.ID].Products, ProductView{
			Name:  product.Name,
			Total: cp.TotalQuantity,
			Items: productItems,
		})
	}

	// Convert map to slice for easier template usage
	var cityWithProductsList []CityWithProductsView
	for _, city := range cityMap {
		cityWithProductsList = append(cityWithProductsList, *city)
	}

	return cityWithProductsList, nil
}
