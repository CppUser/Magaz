package handler

import (
	"Magaz/internal/config"
	"Magaz/pkg/bot/telegram"
)

// TODO: Needs to move to more global space so everyone who need can access it
// Mb call Repository ?
type Handler struct {
	Api *config.APIConfig
	Bot *telegram.Bot
}

func NewHandler(api *config.APIConfig, bot *telegram.Bot) *Handler {
	return &Handler{
		Api: api,
		Bot: bot,
	}
}
