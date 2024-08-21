package telegram

import (
	models2 "Magaz/internal/storage/models"
	"Magaz/pkg/bot/telegram/handlers"
	"Magaz/pkg/utils/convert"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserStatus struct {
	ID    int64  `json:"id"`
	State string `json:"state"`
}

// TODO: refactor code
func CheckUserExists(b *Bot, update telego.Update) (*UserStatus, error) {
	ctx := context.Background()
	userKey := fmt.Sprintf("id:%d", update.Message.From.ID)

	// Check if user exists in cache
	userData, err := b.Cache.Get(ctx, userKey).Result()
	if err != nil {
		b.Logger.Info("User not found in cache", zap.Int64("telegram_id", update.Message.From.ID))

		b.Logger.Info("Checking if user exists in DB", zap.Int64("telegram_id", update.Message.From.ID))
		var user models2.User
		result := b.DB.Select("telegram_id", "status").First(&user, "telegram_id = ?", update.Message.From.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			b.Logger.Info("User not found in DB, creating new user in DB and in cache ", zap.Int64("telegram_id", update.Message.From.ID))

			// Create new user in DB
			user = models2.User{
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

func FetchCitiesFromDB(db *gorm.DB) ([]handlers.TempMarkup, error) {
	var cities []models2.City
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
