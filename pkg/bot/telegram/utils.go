package telegram

import (
	"Magaz/internal/storage/crud"
	"Magaz/internal/storage/models"
	"Magaz/pkg/bot/telegram/handlers"
	"Magaz/pkg/utils/convert"
	"Magaz/pkg/utils/state/fsm"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"time"
)

//TODO: Refactor entire file

type UserStatus struct {
	ID    int64  `json:"id"`
	State string `json:"state"`
}

// TODO:CheckUserExists refactor code
func CheckUserExists(b *Bot, update telego.Update) (*UserStatus, error) {
	ctx := context.Background()
	userKey := fmt.Sprintf("id:%d", update.Message.From.ID)

	// Check if user exists in cache
	userData, err := b.Cache.Get(ctx, userKey).Result()
	if err != nil {
		b.Logger.Info("User not found in cache", zap.Int64("telegram_id", update.Message.From.ID))

		b.Logger.Info("Checking if user exists in DB", zap.Int64("telegram_id", update.Message.From.ID))
		var user models.User
		result := b.DB.Select("telegram_id", "status").First(&user, "telegram_id = ?", update.Message.From.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			b.Logger.Info("User not found in DB, creating new user in DB and in cache ", zap.Int64("telegram_id", update.Message.From.ID))

			// Create new user in DB
			user = models.User{
				TelegramID: update.Message.From.ID,
				Username:   update.Message.From.Username,
				FirstName:  update.Message.From.FirstName,
				LastName:   update.Message.From.LastName,
				Language:   update.Message.From.LanguageCode,
				Status:     "",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := b.DB.Create(&user).Error; err != nil {
				b.Logger.Error("Failed to create new user in DB", zap.String("error", err.Error()))
				return nil, err
			}

			// Set user in cache
			userStatus := UserStatus{
				ID:    user.TelegramID,
				State: user.Status,
			}
			// Set user in cache
			cachedUser, err := convert.ToJSON(userStatus)
			if err == nil {
				err = b.Cache.Set(ctx, userKey, cachedUser, 0).Err()
				if err != nil {
					b.Logger.Error("Failed to add user to cache", zap.Int64("telegram_id", update.Message.From.ID))
				} else {
					b.Logger.Info("User added to cache", zap.Int64("telegram_id", update.Message.From.ID))
				}

			}
			return &userStatus, nil
		} else {
			b.Logger.Info("User found in DB", zap.Int64("telegram_id", update.Message.From.ID))

			// Set user in cache
			userStatus := UserStatus{
				ID:    user.TelegramID,
				State: user.Status,
			}
			cachedUser, err := convert.ToJSON(userStatus)
			if err == nil {
				err = b.Cache.Set(ctx, userKey, cachedUser, 0).Err()
				if err != nil {
					b.Logger.Error("Failed to add user to cache", zap.Int64("telegram_id", update.Message.From.ID))
				} else {
					b.Logger.Info("User added to cache", zap.Int64("telegram_id", update.Message.From.ID))
				}
			}
		}
	}

	b.Logger.Info("User found in cache", zap.Int64("telegram_id", update.Message.From.ID))

	var userStatus UserStatus
	//TODO: add to converter pkg func to convert from JSON
	if err := json.Unmarshal([]byte(userData), &userStatus); err != nil {
		b.Logger.Error("Failed to unmarshal user data", zap.String("error", err.Error()))
		return nil, err
	}
	return &userStatus, nil
}

// TODO:Refactor code
func DefineRules(bot *telego.Bot, db *gorm.DB, event string, ctx string) []fsm.Rule {
	rules := []fsm.Rule{
		{
			Event:      fsm.Event(event),
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					// Fetch the city name from the context (assuming it's passed in the context)
					cityName := context[ctx].(string)

					// Generate product markup for the city
					markup, err := GenerateProductMarkup(db, cityName)
					if err != nil {
						return err
					}

					// Use the generated markup in the action
					return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующий вас товар:", markup)(context)
				},
			},
		},
		// Add more rules here...
	}

	return rules
}

func FetchCitiesFromDB(db *gorm.DB) ([]handlers.TempMarkup, error) {
	var cities []models.City
	err := db.Find(&cities).Error
	if err != nil {
		return nil, err
	}

	var markup []handlers.TempMarkup
	for _, city := range cities {
		markup = append(markup, handlers.TempMarkup{
			Text:         city.Name,
			CallbackData: "city:" + city.Name,
		})
	}

	return markup, nil
}

// GenerateProductMarkup generates the inline keyboard markup for products based on the city name.
func GenerateProductMarkup(db *gorm.DB, cityName string) ([]handlers.TempMarkup, error) {
	// Fetch products associated with the city by name
	products, err := crud.GetCityProducts(db, cityName)
	if err != nil {
		return nil, err
	}

	// Generate markup
	markup := make([]handlers.TempMarkup, len(products))
	for i, product := range products {
		markup[i] = handlers.TempMarkup{
			Text:         product.Name,
			CallbackData: "product:" + product.Name, // Assuming you want to send product name as callback data
		}
	}

	return markup, nil
}

func GenerateProductPriceMarkup(db *gorm.DB, productName string) ([]handlers.TempMarkup, error) {
	log.Println("Generating markup for product:", productName)

	productPrices, err := crud.GetProductPricesByName(db, productName)
	if err != nil {
		return nil, err
	}

	markup := make([]handlers.TempMarkup, len(productPrices))
	for i, productPrice := range productPrices {
		fmt.Println(fmt.Sprintf("Количество %v - Цена %.2f", productPrice.Quantity, productPrice.Price))

		quantityText := fmt.Sprintf("Количество %v - Цена %.2f", productPrice.Quantity, productPrice.Price)
		markup[i] = handlers.TempMarkup{
			Text:         quantityText,
			CallbackData: "quantity:" + fmt.Sprintf("%v", productPrice.Quantity),
		}
	}
	log.Println("Generated markup:", markup)
	return markup, nil
}

func StoreUserChoice(rd *redis.Client, userID int64, step string, choice string) {
	key := fmt.Sprintf("user:%d", userID)
	_, err := rd.HSet(context.Background(), key, step, choice).Result()
	if err != nil {
		log.Println("Failed to store user choice in Redis", err)
	}
}

func GetUserChoices(rd *redis.Client, userID int64) (map[string]string, error) {
	key := fmt.Sprintf("user:%d", userID)
	return rd.HGetAll(context.Background(), key).Result()
}
