package repository

import (
	crud "Magaz/backend/internal/storage"
	"Magaz/backend/internal/storage/models"
	"fmt"
	"gorm.io/gorm"
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

func GetUnreleasedOrders(db *gorm.DB) ([]OrderView, error) {
	var orders []models.Order
	err := db.Preload("User").Preload("City").Preload("Product").Preload("AddrToRelease").Where("released = ?", false).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	var orderViews []OrderView
	for _, order := range orders {
		orderView := OrderView{
			ID:          order.ID,
			ProductName: order.Product.Name,
			CityName:    order.City.Name,
			Quantity:    order.Quantity,
			Due:         order.Due,
			CreatedAt:   order.CreatedAt,
		}

		// Safely handle nil User association
		if order.User != nil {
			orderView.Client = UserView{
				ID:        order.User.ID,
				ChatID:    order.User.ChatID,
				Username:  order.User.Username,
				FirstName: order.User.FirstName,
				LastName:  order.User.LastName,
			}
		} else {
			orderView.Client = UserView{} // Or handle it as per your application's requirement
		}

		// Safely handle nil AddrToRelease association
		if order.AddrToRelease != nil {
			orderView.Address = AddressView{
				ID:          order.AddrToRelease.ID,
				City:        order.AddrToRelease.City.Name,
				Description: order.AddrToRelease.Description,
				Image:       order.AddrToRelease.Image,
				AddedAt:     order.AddrToRelease.AddedAt,
				AddedBy:     order.AddrToRelease.AddedBy.Username,
			}
		} else {
			orderView.Address = AddressView{} // Or handle it as per your application's requirement
		}

		// Populate PaymentMethod based on the type (same logic as before)
		if order.PaymentMethodType == "card" {

			card, err := crud.Get[models.Card](db, order.PaymentMethodID)
			if err != nil {
				return nil, fmt.Errorf("failed to get card payment method: %w", err)
			}

			orderView.PaymentMethod = PaymentView{
				PaymentCategory: "Карта & СБП",
				CardPayment: CardView{
					BankName:   card.BankName,
					BankUrl:    card.BankURL,
					CardNumber: card.CardNumber,
					FirstName:  card.FirstName,
					LastName:   card.LastName,
					UserName:   card.UserID,
					Password:   card.Password,
					QuickPay:   card.QuickPay,
				},
			}
		} else if order.PaymentMethodType == "crypto" {

			cr, err := crud.Get[models.Crypto](db, order.PaymentMethodID)
			if err != nil {
				return nil, fmt.Errorf("failed to get card payment method: %w", err)
			}

			orderView.PaymentMethod = PaymentView{
				PaymentCategory: "Crypto",
				CryptoPayment: CryptoView{
					WalletName:    cr.WalletName,
					WalletAddress: cr.WalletAddr,
					WalletURL:     "TODO:Add to crypto model",
					UserName:      cr.UserID,
					Password:      cr.Password,
				},
			}
		}

		orderViews = append(orderViews, orderView)
	}

	return orderViews, nil
}
