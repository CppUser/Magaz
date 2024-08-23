package crud

import (
	"Magaz/internal/storage/models"
	"gorm.io/gorm"
	"log"
)

func GetProductPricesByName(db *gorm.DB, productName string) ([]models.ProductPrice, error) {
	var product models.Product

	// Fetch the product by its name
	if err := db.Where("name = ?", productName).First(&product).Error; err != nil {
		log.Println("Error fetching product:", err)
		return nil, err
	}
	log.Println("Fetched product:", product)

	// Fetch the prices and quantities for the product
	var productPrices []models.ProductPrice
	if err := db.Where("product_id = ?", product.ID).Find(&productPrices).Error; err != nil {
		log.Println("Error fetching product prices:", err)
		return nil, err
	}

	log.Println("Fetched product prices:", productPrices)
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
