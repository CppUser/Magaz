package repository

import (
	"Magaz/backend/internal/storage/models"

	"gorm.io/gorm"
)

//TODO: This is temporary file to fetch data from db
//TODO: Refactor this file . storage crud exist for it

// ProductItem represents a specific quantity and price of a product.
type ProductItem struct {
	QuantityID uint
	Quantity   float32
	Price      float32
	AddrCnt    int64 //How many addresses of that product
}

// ProductView is used for displaying products in the UI.
type ProductView struct {
	Name      string
	Total     float32 // This represents the total available quantity
	ProductID uint
	Items     []ProductItem
}

// CityWithProductsView represents a city and its associated products for display.
type CityWithProductsView struct {
	CityID   uint
	City     string
	Products []ProductView
}

func FetchCityProcdducts(db *gorm.DB) ([]CityWithProductsView, error) {
	var cityProducts []models.CityProduct

	//
	err := db.Preload("City").
		Preload("Product").
		Preload("QtnPrices").
		Preload("QtnPrices.Address").
		Find(&cityProducts).Error
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
				CityID:   city.ID,
				City:     city.Name,
				Products: []ProductView{},
			}
		}

		// Convert ProductPrices to ProductItems
		var productItems []ProductItem
		for _, pp := range cp.QtnPrices {

			var addrCount int64
			err := db.Model(&models.Address{}).
				Where("qtn_price_id = ? AND released = ?", pp.ID, false).
				Count(&addrCount).Error
			if err != nil {
				return nil, err
			}
			productItems = append(productItems, ProductItem{
				QuantityID: pp.ID,
				Quantity:   pp.Quantity,
				Price:      pp.Price,
				AddrCnt:    addrCount,
			})
		}

		// Add product under the city
		cityMap[city.ID].Products = append(cityMap[city.ID].Products, ProductView{
			ProductID: product.ID,
			Name:      product.Name,
			Total:     cp.TotalQuantity,
			Items:     productItems,
		})
	}

	// Convert map to slice for easier template usage
	var cityWithProductsList []CityWithProductsView
	for _, city := range cityMap {
		cityWithProductsList = append(cityWithProductsList, *city)
	}

	return cityWithProductsList, nil
}
