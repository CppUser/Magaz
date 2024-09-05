package models

import "time"

type Order struct {
	ID                uint `gorm:"primaryKey;autoIncrement"`
	UserID            int64
	User              *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	CityID            uint    `gorm:"not null"`
	City              City    `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	ProductID         uint    `gorm:"not null"`
	Product           Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Quantity          float32 `gorm:"not null"`
	Due               uint    `gorm:"not null"`
	PaymentMethodType string  `gorm:"size:100;not null"` // e.g., "Card", "Crypto"
	PaymentMethodID   uint    `gorm:"not null"`          // ID of the associated payment method
	PaymentConfImg    string
	Released          bool `gorm:"default:false"`
	ReleasedByID      *uint
	ReleasedBy        *Employee `gorm:"foreignKey:ReleasedByID;constraint:OnDelete:CASCADE;"`
	ReleaseTime       time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
