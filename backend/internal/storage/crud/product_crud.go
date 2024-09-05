package crud

import (
	"Magaz/backend/internal/storage/models"
	"fmt"
	"gorm.io/gorm"
)

func GetProductPricesByName(db *gorm.DB, productName string) ([]models.QtnPrice, error) {
	var product models.Product

	// Fetch the product by its name
	if err := db.Where("name = ?", productName).First(&product).Error; err != nil {
		return nil, err
	}

	// Fetch the prices and quantities for the product
	var productPrices []models.QtnPrice
	if err := db.Where("city_product_id = ?", product.ID).Find(&productPrices).Error; err != nil {
		return nil, err
	}

	return productPrices, nil
}

// GetCityProducts retrieves all products associated with a city by its name.
func GetCityProducts(db *gorm.DB, cityName string) ([]models.Product, error) {
	var city models.City

	// First, find the city by its name
	if err := db.Where("name = ?", cityName).First(&city).Error; err != nil {
		return nil, err
	}

	// Then, find the products associated with the city
	var cityProducts []models.CityProduct
	if err := db.Preload("Product").Where("city_id = ?", city.ID).Find(&cityProducts).Error; err != nil {
		return nil, err
	}

	// Extract the products from cityProducts
	products := make([]models.Product, len(cityProducts))
	for i, cp := range cityProducts {
		products[i] = cp.Product
	}

	return products, nil
}

func GetProductIDByCityAndProductName(db *gorm.DB, cityName string, productName string) (uint, error) {
	var city models.City

	// Find the city ID by name
	if err := db.Where("name = ?", cityName).First(&city).Error; err != nil {
		return 0, fmt.Errorf("city not found: %w", err)
	} // retrieve 2

	// Then, find the products associated with the city
	var cityProducts []models.CityProduct
	if err := db.Preload("Product").Where("city_id = ?", city.ID).Find(&cityProducts).Error; err != nil {
		return 0, err
	} //retrieve all products in that city 2

	for _, p := range cityProducts {
		if p.Product.Name == productName {
			return p.ProductID, nil
		}

	}
	// Return the ProductID from CityProduct
	return 0, fmt.Errorf("product with name '%s' not found in city '%s'", productName, cityName)
}
