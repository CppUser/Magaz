package models

import "time"

type Employee struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"size:100"`
	LastName  string `gorm:"size:100"`
	Username  string `gorm:"size:100"`
	Password  string `gorm:"size:100"`
	Email     string `gorm:"size:100"`
	Phone     string `gorm:"size:20"`
	Role      string `gorm:"size:100"`
	Active    bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
