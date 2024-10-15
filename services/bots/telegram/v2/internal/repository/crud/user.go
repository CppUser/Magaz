package crud

import (
	"context"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"tg/internal/bot"
	"tg/internal/repository/models"
	"time"
)

// TODO:Refactor CheckUserDetails
func CheckUserDetails(b *bot.BotService, update telego.Update) (*models.CachedUserState, error) {
	ctx := context.Background()
	userKey := fmt.Sprintf("user:%d", update.Message.From.ID)

	// Check if the user hash exists in Redis
	var exists, err = b.Rdb.Exists(ctx, userKey).Result()
	if err != nil {
		b.Logger.Error("Error checking user in cache", zap.String("error", err.Error()))
		return nil, err
	}

	if exists > 0 {
		userData, err := b.Rdb.HGetAll(ctx, userKey).Result()
		if err != nil {
			b.Logger.Error("Failed to retrieve user data from cache", zap.String("error", err.Error()))
			return nil, err
		}

		// Populate UserCache from hash fields
		userStatus := models.CachedUserState{
			ID:          update.Message.From.ID,
			OrderStatus: userData["order_status"],
			State:       userData["state"],
		}
		return &userStatus, nil
	}

	var user models.User
	err = b.DB.First(&user, "id = ?", update.Message.From.ID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = models.User{
			ID:        update.Message.From.ID,
			ChatID:    update.Message.GetChat().ID,
			Username:  update.Message.From.Username,
			FirstName: update.Message.From.FirstName,
			LastName:  update.Message.From.LastName,
			Language:  update.Message.From.LanguageCode,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := b.DB.Create(&user).Error; err != nil {
			b.Logger.Error("Failed to create new user in DB", zap.String("error", err.Error()))
			return nil, err
		}

		usr := models.CachedUserState{
			ID:          user.ID,
			OrderStatus: "initial",
			State:       "/start",
		}
		cacheUserAsHash(ctx, b, userKey, usr)
		return &usr, nil
	} else if err != nil {
		b.Logger.Error("Failed to retrieve user from DB", zap.String("error", err.Error()))
		return nil, err
	}

	// User found in DB, check the latest order
	var latestOrder models.Order
	err = b.DB.Where("user_id = ?", user.ID).Order("created_at DESC").First(&latestOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User has no orders
			b.Logger.Info("No orders found for user", zap.Int64("telegram_id", update.Message.From.ID))
			// Cache the user details with "no orders" status
			userStatus := models.CachedUserState{
				ID:          user.ID,
				OrderStatus: "initial",
			}
			cacheUserAsHash(ctx, b, userKey, userStatus)
			return &userStatus, nil
		} else {
			b.Logger.Error("Failed to retrieve user's latest order", zap.String("error", err.Error()))
			return nil, err
		}
	}

	userStatus := models.CachedUserState{
		ID: user.ID,
	}

	//TODO: Might run into issue when operator released on his end ,
	//client may ask for status it will lead to new /start state
	if latestOrder.Released {
		userStatus.OrderStatus = "released"
		userStatus.State = "/start"
	} else {
		userStatus.OrderStatus = "processing"
	}

	cacheUserAsHash(ctx, b, userKey, userStatus)
	return &userStatus, nil

}

func cacheUserAsHash(ctx context.Context, b *bot.BotService, userKey string, usr models.CachedUserState) {
	err := b.Rdb.HMSet(ctx, userKey, map[string]interface{}{
		"order_status": usr.OrderStatus,
		"state":        usr.State,
	}).Err()

	if err != nil {
		b.Logger.Error("Failed to add user to cache", zap.Int64("telegram_id", usr.ID))
		return // Early return if there was an error setting the hash
	}

	// Set the lifetime (TTL) for the hash key (e.g., 1 hour)
	if err := b.Rdb.Expire(ctx, userKey, time.Hour-55).Err(); err != nil {
		b.Logger.Error("Failed to set expiration for user hash", zap.Int64("telegram_id", usr.ID))
	} else {
		b.Logger.Info("User added to cache with a 1-hour expiration", zap.Int64("telegram_id", usr.ID))
	}
}
