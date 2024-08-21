package handler

import (
	"Magaz/internal/config"
	"Magaz/pkg/bot/telegram"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TODO: Needs to move to more global space so everyone who need can access it
// Mb call Repository ?
type Handler struct {
	Api    *config.APIConfig
	Logger *zap.Logger
	Bot    *telegram.Bot
	Redis  *redis.Client
	DB     *gorm.DB
}

//func (h *Handler) NewHandler() *Handler {
//	return &Handler{
//		Api:    api,
//		Logger: log,
//		Bot:    bot,
//		Redis:  rd,
//		DB:     db,
//	}
//}
