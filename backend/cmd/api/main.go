package main

import (
	"Magaz/backend/internal/config"
	"Magaz/backend/internal/handler"
	"Magaz/backend/internal/router"
	"Magaz/backend/internal/storage/models"
	"Magaz/backend/internal/utils"
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
	//bot := telegram.Bot{
	//	Config:           &cfg.Bot,
	//	Logger:           zaplog,
	//	UpdateChanBuffer: 128, // Buffer size is 128 default
	//	Cache:            rdb,
	//	DB:               db,
	//}
	//bot.InitBot()

	tempalteCache, err := handler.CreateTemplateCache(cfg.CacheDir.Layouts, cfg.CacheDir.Pages)
	if err != nil {
		zaplog.Fatal("Failed to create template cache", zap.Error(err))
	}

	h := handler.Handler{
		Api:    cfg,
		Logger: zaplog,
		//Bot:       &bot,
		Redis:     rdb,
		DB:        db,
		TmplCache: tempalteCache,
		Session:   store,
	}
	//handler.NewHandler(h)

	rh := router.SetupRouter(&h)

	//go bot.ReceiveUpdates() //TODO: no the best approach find other way to handle updates

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
		//&models.User{},
		&models.City{},
		&models.Product{},
		&models.ProductPrice{},
		&models.CityProduct{},
		&models.Item{},
	)
}

func PopulateData(db *gorm.DB) error {
	// Sample Cities
	newYork := models.City{Name: "New York"}
	chicago := models.City{Name: "Chicago"}

	// Insert Cities into the database
	db.Create(&newYork)
	db.Create(&chicago)

	// Sample Products
	apples := models.Product{
		Name:        "Apples",
		Description: "Fresh apples",
		Image:       "apples.jpg",
	}
	bananas := models.Product{
		Name:        "Bananas",
		Description: "Ripe bananas",
		Image:       "bananas.jpg",
	}

	// Insert Products into the database
	db.Create(&apples)
	db.Create(&bananas)

	// Sample CityProducts with Prices and Quantities
	nyApples := models.CityProduct{
		CityID:        newYork.ID,
		ProductID:     apples.ID,
		TotalQuantity: 300, // 300kg of apples in New York
		ProductPrices: []models.ProductPrice{
			{Quantity: 1, Price: 10}, // $10 per kg
			{Quantity: 5, Price: 45}, // $45 per 5kg
		},
	}
	nyBananas := models.CityProduct{
		CityID:        newYork.ID,
		ProductID:     bananas.ID,
		TotalQuantity: 200, // 200kg of bananas in New York
		ProductPrices: []models.ProductPrice{
			{Quantity: 1, Price: 3},  // $3 per kg
			{Quantity: 5, Price: 14}, // $14 per 5kg
		},
	}
	chiApples := models.CityProduct{
		CityID:        chicago.ID,
		ProductID:     apples.ID,
		TotalQuantity: 100, // 100kg of apples in Chicago
		ProductPrices: []models.ProductPrice{
			{Quantity: 1, Price: 22},  // $22 per kg
			{Quantity: 5, Price: 100}, // $100 per 5kg
		},
	}

	// Insert CityProducts with associated ProductPrices into the database
	db.Create(&nyApples)
	db.Create(&nyBananas)
	db.Create(&chiApples)

	return nil
}
