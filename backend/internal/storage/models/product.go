package models

// Product represents a type of product that can be sold in various cities.
type Product struct {
	ID           uint          `gorm:"primaryKey"`
	Name         string        `gorm:"size:100;not null"`
	Description  string        `gorm:"size:255"`
	Image        string        `gorm:"size:255"`
	CityProducts []CityProduct `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Items        []Item        `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

// CityProduct represents the availability and pricing details of a product in a specific city.
type CityProduct struct {
	ID            uint           `gorm:"primaryKey"`
	CityID        uint           `gorm:"index;not null"`
	ProductID     uint           `gorm:"index;not null"`
	TotalQuantity float32        `gorm:"not null"` // Available quantity of the product in the city
	Product       Product        `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	City          City           `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	ProductPrices []ProductPrice `gorm:"foreignKey:CityProductID;constraint:OnDelete:CASCADE;"`
}

// ProductPrice represents the price of a product based on quantity in a specific city.
type ProductPrice struct {
	ID            uint        `gorm:"primaryKey"`
	CityProductID uint        `gorm:"index;not null"`
	Quantity      float32     `gorm:"not null"` // Quantity for which this price applies (e.g., per kg)
	Price         float32     `gorm:"not null"` // Price for the given quantity
	CityProduct   CityProduct `gorm:"foreignKey:CityProductID;constraint:OnDelete:CASCADE;"`
}
