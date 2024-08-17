package main

import (
	"Magaz/internal/config"
	"Magaz/internal/handler"
	"Magaz/internal/router"
	"Magaz/pkg/bot/telegram"
	"Magaz/pkg/client/redis"
	"Magaz/pkg/utils/logger"
	"go.uber.org/zap"
	"log"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int    `json:"message_id"`
		Text      string `json:"text"`
		Chat      struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load API configs: %v", err)
	}

	cfg.Logger, _ = logger.InitLogger(cfg.Env)
	cfg.Bot.Logger = cfg.Logger //TODO: currently sharing logger . find other way to share logger

	rdb, rdberr := redis.InitRedisClient(&cfg.Redis) //TODO: assign to var
	if rdberr != nil {
		cfg.Logger.Fatal("Failed to connect to Redis", zap.String("error", rdberr.Error()))
	}

	//TODO: Initialize Sessions
	//TODO: Initialize DB

	bot := telegram.Bot{
		Config:           &cfg.Bot,
		UpdateChanBuffer: 128, // Buffer size is 128 default
		Cache:            rdb,
	}
	bot.InitBot()

	h := handler.NewHandler(cfg, &bot)
	rh := router.SetupRouter(h)

	go bot.ReceiveUpdates() //TODO: no the best approach find other way to handle updates

	//TODO: build host and port from config and call single function
	err = rh.Run(cfg.Server.Host + ":" + cfg.Server.Port)
	if err != nil {
		cfg.Logger.Fatal("Failed to start server", zap.String("error", err.Error()))
	}

}
