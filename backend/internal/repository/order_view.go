package repository

import (
	"time"
)

type OrderView struct {
	ID            uint        `json:"id"`
	ProductName   string      `json:"product_name"`
	CityName      string      `json:"city_name"`
	Quantity      float32     `json:"quantity"`
	Due           uint        `json:"due"`
	CreatedAt     time.Time   `json:"created_at"`
	Client        UserView    `json:"user_view"`
	PaymentMethod PaymentView `json:"payment_method"`
	Address       AddressView `json:"address"`
}
