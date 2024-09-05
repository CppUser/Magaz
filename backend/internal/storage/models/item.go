package models

import "time"

type Address struct {
	ID          uint      `gorm:"primaryKey"`
	CityID      uint      `gorm:"index;not null"`
	QtnPriceID  uint      `gorm:"index;not null"` // Quantity of the product at this address.
	ProductID   uint      `gorm:"index;not null"`
	Description string    `gorm:"size:255"` // Details about the storage location or conditions.
	Image       string    `gorm:"size:255"` // Image of the location, product storage, etc.
	AddedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EmployeeID  uint      `gorm:"not null"` // The employee who added this record.
	AddedBy     Employee  `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;"`
	ReleaseDate time.Time `gorm:""`
	//UserTelegramID uint      `gorm:"not null"`
	//ReleasedTo     User      `gorm:"foreignKey:UserTelegramID;constraint:OnDelete:CASCADE;"` //TODO: can be released to not only user
	Product  Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	City     City     `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	QtnPrice QtnPrice `gorm:"foreignKey:QtnPriceID;constraint:OnDelete:CASCADE;"`
}
