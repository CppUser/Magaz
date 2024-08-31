package models

import "time"

type Item struct {
	ID          uint      `gorm:"primaryKey"`
	ProductID   uint      `gorm:"index;not null"`
	Quantity    float32   `gorm:"not null"`
	Description string    `gorm:"size:255"`
	Image       string    `gorm:"size:255"`
	AddedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	//AddedBy     Employee  `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;"`
	//ReleaseDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	//ReleasedTo  User      `gorm:"foreignKey:UserTelegramID;constraint:OnDelete:CASCADE;"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}
