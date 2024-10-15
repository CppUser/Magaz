package models

import "time"

type User struct {
	ID            int64  `gorm:"primaryKey;uniqueIndex"`
	Type          string `gorm:""`
	ChatID        int64  `gorm:""`
	Username      string `gorm:"size:100"`
	FirstName     string `gorm:"size:100"`
	LastName      string `gorm:"size:100"`
	Language      string `gorm:"size:10"`
	PhoneNumber   string `gorm:"size:20"`
	AccessLevel   uint8  `gorm:"default:1"`
	Blocked       bool   `gorm:"default:false"`
	BlockedTime   time.Time
	BlockedReason string `gorm:"size:255"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
