package models

import "time"

type Address struct {
	ID             uint      `gorm:"primaryKey"`
	CityID         uint      `gorm:"index;not null"`
	QtnPriceID     uint      `gorm:"index;not null"` // Quantity of the product at this address.
	ProductID      uint      `gorm:"index;not null"`
	Description    string    `gorm:"size:255"` // Details about the storage location or conditions.
	Image          string    `gorm:"size:255"` // Image of the location, product storage, etc.
	AddedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EmployeeID     uint      `gorm:""` // The employee who added this record.
	AddedBy        Employee  `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;"`
	Assigned       bool      `gorm:"default:false"`
	AssignedUserID *int64    `gorm:""`
	AssignedTo     User      `gorm:"foreignKey:AssignedUserID;constraint:OnDelete:CASCADE;"`
	Released       bool      `gorm:"default:false"`
	ReleaseDate    time.Time `gorm:""`
	ReleasedTo     string    `gorm:""`
	Product        Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	City           City      `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"`
	QtnPrice       QtnPrice  `gorm:"foreignKey:QtnPriceID;constraint:OnDelete:CASCADE;"`
}
