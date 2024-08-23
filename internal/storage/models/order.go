package models

import "time"

type Order struct {
	ID                uint     `gorm:"primaryKey"`
	User              User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Ciy               City     `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	Area              Area     `gorm:"foreignKey:AreaID;constraint:OnDelete:CASCADE;"`
	Product           Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Quantity          uint     `gorm:"not null"`
	Due               uint     `gorm:"not null"`
	PaymentMethodType string   `gorm:"size:100;not null"` // e.g., "Card", "Crypto"
	PaymentMethodID   uint     `gorm:"not null"`          // ID of the associated payment method
	Released          bool     `gorm:"default:false"`
	ReleasedBy        Employee `gorm:"constraint:OnDelete:CASCADE;"`
	ReleaseTime       time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
