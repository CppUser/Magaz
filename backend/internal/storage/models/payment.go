package models

import "time"

type Card struct {
	ID         uint   `gorm:"primaryKey"`
	BankName   string `gorm:"size:100"`
	BankURL    string `gorm:"size:100"`
	UserID     string `gorm:"size:100"`
	Password   string `gorm:"size:100"`
	CardNumber string `gorm:"size:100"`
	QuickPay   string `gorm:"size:100"`
	FirstName  string `gorm:"size:100"`
	LastName   string `gorm:"size:100"`
	ExpireDate string `gorm:"size:100"`
	CVV        string `gorm:"size:100"`
	CardType   string `gorm:"size:100"`
	Balance    int    `gorm:"size:100"`
	UsedTimes  uint   `gorm:"default:0"`
	Active     bool   `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Crypto struct {
	ID            uint   `gorm:"primaryKey"`
	WalletName    string `gorm:"size:100;not null"`
	WalletAddr    string `gorm:"size:100;not null"`
	TransactionID string `gorm:"size:100;not null"`
	UserID        string `gorm:"size:100;not null"`
	Password      string `gorm:"size:100;not null"`
	PrivateKey    string `gorm:"size:100;not null"`
	PublicKey     string `gorm:"size:100;not null"`
	Active        bool   `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
