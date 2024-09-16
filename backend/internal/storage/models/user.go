package models

import "time"

type User struct {
	ID          int64  `gorm:"primaryKey;uniqueIndex"`
	ChatID      int64  `gorm:""`
	Username    string `gorm:"size:100"`
	FirstName   string `gorm:"size:100"`
	LastName    string `gorm:"size:100"`
	Language    string `gorm:"size:10"`
	PhoneNumber string `gorm:"size:20"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
