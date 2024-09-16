package telegram

import (
	"Magaz/backend/internal/storage/crud"
	"Magaz/backend/internal/storage/models"
	"Magaz/backend/pkg/bot/telegram/handlers"
	"Magaz/backend/pkg/utils/convert"
	"Magaz/backend/pkg/utils/state/fsm"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mymmrac/telego"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//TODO: Refactor entire file

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
		quantityText := fmt.Sprintf("%v - %.2f", productPrice.Quantity, productPrice.Price)
		markup[i] = handlers.TempMarkup{
			Text:         quantityText,
			CallbackData: "quantity:" + fmt.Sprintf("%v", productPrice.Quantity),
		}
	}
	log.Println("Generated markup:", markup)
	return markup, nil
}

///////////////////REDIS OPERATIONS////////////////////////////////////////////////

type UserCache struct {
	ID          int64  `json:"id"`
	OrderStatus string `json:"status"`
	State       string `json:"state"`
}

// CheckUserDetails checks latest user details
func CheckUserDetails(b *Bot, update telego.Update) (*UserCache, error) {
	ctx := context.Background()
	userKey := fmt.Sprintf("user:%d", update.Message.From.ID)

	// Check if the user hash exists in Redis
	exists, err := b.Cache.Exists(ctx, userKey).Result()
	if err != nil {
		b.Logger.Error("Error checking user in cache", zap.String("error", err.Error()))
		return nil, err
	}

	if exists > 0 {
		userData, err := b.Cache.HGetAll(ctx, userKey).Result()
		if err != nil {
			b.Logger.Error("Failed to retrieve user data from cache", zap.String("error", err.Error()))
			return nil, err
		}

		// Populate UserCache from hash fields
		userStatus := UserCache{
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

		userStatus := UserCache{
			ID:          user.ID,
			OrderStatus: "initial",
			State:       "/start",
		}
		cacheUserAsHash(ctx, b, userKey, userStatus)
		return &userStatus, nil
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
			userStatus := UserCache{
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

	userStatus := UserCache{
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

func cacheUserAsStr(ctx context.Context, b *Bot, userKey string, userStatus UserCache) {
	cachedUser, err := convert.ToJSON(userStatus)
	if err == nil {
		err = b.Cache.Set(ctx, userKey, cachedUser, 0).Err()
		if err != nil {
			b.Logger.Error("Failed to add user to cache", zap.Int64("telegram_id", userStatus.ID))
		} else {
			b.Logger.Info("User added to cache", zap.Int64("telegram_id", userStatus.ID))
		}
	} else {
		b.Logger.Error("Failed to convert user status to JSON", zap.String("error", err.Error()))
	}
}

func cacheUserAsHash(ctx context.Context, b *Bot, userKey string, userStatus UserCache) {
	err := b.Cache.HMSet(ctx, userKey, map[string]interface{}{
		"order_status": userStatus.OrderStatus,
		"state":        userStatus.State,
	}).Err()

	if err != nil {
		b.Logger.Error("Failed to add user to cache", zap.Int64("telegram_id", userStatus.ID))
		return // Early return if there was an error setting the hash
	}

	// Set the lifetime (TTL) for the hash key (e.g., 1 hour)
	if err := b.Cache.Expire(ctx, userKey, time.Hour-55).Err(); err != nil {
		b.Logger.Error("Failed to set expiration for user hash", zap.Int64("telegram_id", userStatus.ID))
	} else {
		b.Logger.Info("User added to cache with a 1-hour expiration", zap.Int64("telegram_id", userStatus.ID))
	}
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
