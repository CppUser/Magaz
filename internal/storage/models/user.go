package models

import "time"

type User struct {
	ID          uint   `gorm:"primaryKey"`
	TelegramID  int64  `gorm:"uniqueIndex"`
	Username    string `gorm:"size:100"`
	FirstName   string `gorm:"size:100"`
	LastName    string `gorm:"size:100"`
	Language    string `gorm:"size:10"`
	PhoneNumber string `gorm:"size:20"`
	Status      string `gorm:"size:20"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
