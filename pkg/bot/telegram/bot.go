package telegram

import (
	"Magaz/internal/config"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"
)

type Bot struct {
	API    *telego.Bot
	Config *config.TGBotConfig
	Logger *zap.Logger
}

func InitBotConfig(cfg *config.TGBotConfig, logger *zap.Logger) *Bot {
	return &Bot{
		Config: cfg,
		Logger: logger,
	}

}

// TODO: refactor code move some logic to handlers
// RunBot runs the Telegram bot
func (b *Bot) RunBot() {
	bot, err := telego.NewBot(b.Config.Token)
	if err != nil {
		//TODO: need to handle the error differently without direct call to zap.String
		b.Logger.Fatal("failed create new bot api instance", zap.String("error", err.Error()))
	}
	b.API = bot

	//TODO: refer to SetWebhookParams to setup additional parameters (like certificate, pending updates, etc.)
	_ = b.API.SetWebhook(&telego.SetWebhookParams{
		URL: b.Config.WebhookLink + b.Config.WebhookPath,
	})

	info, _ := b.API.GetWebhookInfo()
	b.Logger.Info("Webhook Info", zap.Any("info", info)) //TODO: in prod it needs to be in JSON format

}
