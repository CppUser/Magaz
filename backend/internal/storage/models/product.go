package models

import "time"

type Product struct {
	ID            uint           `gorm:"primaryKey"`
	Name          string         `gorm:"size:100;not null"`
	Description   string         `gorm:"size:255"`
	Image         string         `gorm:"size:255"`
	ProductPrices []ProductPrice `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	AreaProducts  []AreaProduct  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	CityProducts  []CityProduct  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

type ProductPrice struct {
	ID        uint    `gorm:"primaryKey"`
	ProductID uint    `gorm:"index;not null"`
	Quantity  float32 `gorm:"not null"`
	Price     float32 `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

// AreaProduct represents the availability of a product in a specific area.
type AreaProduct struct {
	ID        uint    `gorm:"primaryKey"`
	AreaID    uint    `gorm:"index;not null"`
	ProductID uint    `gorm:"index;not null"`
	Area      Area    `gorm:"foreignKey:AreaID;constraint:OnDelete:CASCADE;"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

// CityProduct represents the availability of a product in a specific city.
type CityProduct struct {
	ID        uint    `gorm:"primaryKey"`
	CityID    uint    `gorm:"index;not null"`
	ProductID uint    `gorm:"index;not null"`
	City      City    `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

type ProductToRelease struct {
	ID          uint      `gorm:"primaryKey"`
	ProductID   uint      `gorm:"index;not null"`
	Quantity    float32   `gorm:"not null"`
	Description string    `gorm:"size:255"`
	Image       string    `gorm:"size:255"`
	AddedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AddedBy     Employee  `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;"`
	ReleaseDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ReleasedTo  User      `gorm:"foreignKey:UserTelegramID;constraint:OnDelete:CASCADE;"`
}
