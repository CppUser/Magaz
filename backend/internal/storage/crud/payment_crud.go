package crud

import (
	"Magaz/backend/internal/storage/models"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func GetPaymentMethod(db *gorm.DB, paymentMethod string) (interface{}, error) {
	switch paymentMethod {
	case "card":
		var cards []models.Card

		if err := db.Where("balance <= ? AND used_times <= ?", 600000, 4).Find(&cards).Error; err != nil {
			return nil, fmt.Errorf("error retrieving valid cards: %w", err)
		}

		if len(cards) == 0 {
			return nil, fmt.Errorf("no valid cards available")
		}

		// Seed random number generator
		rand.New(rand.NewSource(time.Now().UnixNano()))

		// Select a random card ID from the valid cards
		randomIndex := rand.Intn(len(cards))
		randomCard := cards[randomIndex]

		return randomCard, nil

	case "crypto":
		var cryptos []models.Crypto

		// Query for all active cryptos
		if err := db.Where("active = ?", true).Find(&cryptos).Error; err != nil {
			return nil, fmt.Errorf("error retrieving valid cryptos: %w", err)
		}

		if len(cryptos) == 0 {
			return nil, fmt.Errorf("no valid cryptos available")
		}

		// Seed random number generator
		rand.Seed(time.Now().UnixNano())

		// Select a random crypto from the valid cryptos
		randomIndex := rand.Intn(len(cryptos))
		randomCrypto := cryptos[randomIndex]

		return randomCrypto, nil

	default:
		return 0, fmt.Errorf("unknown payment method: %s", paymentMethod)
	}
}
