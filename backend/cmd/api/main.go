package main

import (
	"Magaz/backend/internal/config"
	"Magaz/backend/internal/handler"
	"Magaz/backend/internal/router"
	"Magaz/backend/internal/storage/models"
	"Magaz/backend/internal/utils"
	"Magaz/backend/pkg/bot/telegram"
	"Magaz/backend/pkg/client/postgres"
	"Magaz/backend/pkg/client/redis"
	"Magaz/backend/pkg/utils/logger"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load API configs: %v", err)
	}

	zaplog, _ := logger.InitLogger(cfg.Env)

	rdb, rdberr := redis.InitRedisClient(&cfg.Redis) //TODO: assign to var
	if rdberr != nil {
		zaplog.Fatal("Failed to connect to Redis", zap.String("error", rdberr.Error()))
	}

	db, dberr := postgres.Connect(cfg.Database)
	if dberr != nil {
		zaplog.Fatal("Failed to connect to DB", zap.String("error", dberr.Error()))

	}

	///////////////////////////////////////////////////////////////////////////////
	// Migrate the database models
	_ = migrateDatabase(db)

	//err = PopulateData(db)
	//if err != nil {
	//	fmt.Println("Failed to populate the database:", err)
	//} else {
	//	fmt.Println("Database populated successfully!")
	//}
	////////////////////////////////////////////////////////////////////////////////
	//if err = db.AutoMigrate(&models.User{}); err != nil {
	//	zaplog.Fatal("Failed to migrate database schema", zap.Error(err))
	//}

	sessionKey := cfg.ScrKey
	if sessionKey == "" {
		sessionKey, err = utils.GenerateRandomKey(32)
		if err != nil {
			zaplog.Fatal("Failed to generate session key", zap.Error(err))
		}
	}
	// Initialize the session store with the retrieved or generated session key
	store := sessions.NewCookieStore([]byte(sessionKey))

	//TODO: passing to handler initialized clients like redis and db . Pass handler instead ?
	bot := telegram.Bot{
		Config:           &cfg.Bot,
		Logger:           zaplog,
		UpdateChanBuffer: 128, // Buffer size is 128 default
		Cache:            rdb,
		DB:               db,
	}
	bot.InitBot()

	tempalteCache, err := handler.CreateTemplateCache(cfg.CacheDir.Layouts, cfg.CacheDir.Pages)
	if err != nil {
		zaplog.Fatal("Failed to create template cache", zap.Error(err))
	}

	h := handler.Handler{
		Api:       cfg,
		Logger:    zaplog,
		Bot:       &bot,
		Redis:     rdb,
		DB:        db,
		TmplCache: tempalteCache,
		Session:   store,
	}
	//handler.NewHandler(h)

	rh := router.SetupRouter(&h)

	go bot.ReceiveUpdates() //TODO: no the best approach find other way to handle updates

	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: rh,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zaplog.Fatal("Failed to listen and serve", zap.String("error", err.Error()))

		}
	}()

	<-quit
	zaplog.Info("Shutting down server...")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		zaplog.Fatal("Server forced to shutdown:", zap.Error(err))
	}

}

// Migrate the database models
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.City{},
		&models.Product{},
		&models.QtnPrice{},
		&models.CityProduct{},
		&models.Address{},
		&models.Order{},
		&models.Card{},
		&models.Crypto{},
	)
}
