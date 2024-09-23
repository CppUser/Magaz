package models

import "time"

// TODO: Refactor should be gorm and json or one of them
type Order struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int64     `gorm:""`
	User              *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	CityID            uint      `gorm:"not null"`
	City              City      `gorm:"foreignKey:CityID;constraint:OnDelete:CASCADE;"` //TODO: must be a pointer to city
	ProductID         uint      `gorm:"not null"`
	Product           Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"` //TODO: must be a pointer to product
	Quantity          float32   `gorm:"not null" json:"quantity"`
	Due               uint      `gorm:"not null" json:"due"`
	PaymentMethodType string    `gorm:"size:100;not null" json:"paymentMethodType"` // e.g., "Card", "Crypto"
	PaymentMethodID   uint      `gorm:"not null"`                                   // ID of the associated payment method
	PaymentConfImg    string    `gorm:""`
	Released          bool      `gorm:"default:false"`
	Declined          bool      `gorm:"default:false"`
	ReleasedByID      *uint     `gorm:"" json:"released_by_id"`
	ReleasedBy        *Employee `gorm:"foreignKey:ReleasedByID;constraint:OnDelete:CASCADE;"`
	ReleaseTime       time.Time `gorm:"" json:"release_time"`
	ReleasedAddrID    *uint     `gorm:""`
	AddrToRelease     *Address  `gorm:"foreignKey:ReleasedAddrID;constraint:OnDelete:CASCADE;"`
	CreatedAt         time.Time `gorm:"" json:"created_at"`
	UpdatedAt         time.Time `gorm:"" json:"updated_at"`
}

type DeclinedOrder struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      uint      `gorm:"not null" json:"order_id"`
	Order        Order     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
	Reason       string    `gorm:"not null" json:"reason"`
	DeclinedAt   time.Time `gorm:"" json:"declined_at"`
	DeclinedByID uint      `gorm:"" json:"declined_by"`
	Employee     Employee  `gorm:"foreignKey:DeclinedByID;constraint:OnDelete:CASCADE;"`
}
