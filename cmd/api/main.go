package main

import (
	"Magaz/internal/config"
	"Magaz/internal/handler"
	"Magaz/internal/router"
	"Magaz/internal/storage/models"
	"Magaz/pkg/bot/telegram"
	"Magaz/pkg/client/postgres"
	"Magaz/pkg/client/redis"
	"Magaz/pkg/utils/logger"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
	//_ = migrateDatabase(db)
	//
	//// Populate the database
	//err = populateDatabase(db)
	//if err != nil {
	//	fmt.Println("Failed to populate the database:", err)
	//} else {
	//	fmt.Println("Database populated successfully!")
	//}
	////////////////////////////////////////////////////////////////////////////////
	//if err = db.AutoMigrate(&models.User{}); err != nil {
	//	zaplog.Fatal("Failed to migrate database schema", zap.Error(err))
	//}

	////TODO: Initialize Sessions

	//TODO: passing to handler initialized clients like redis and db . Pass handler instead ?
	bot := telegram.Bot{
		Config:           &cfg.Bot,
		Logger:           zaplog,
		UpdateChanBuffer: 128, // Buffer size is 128 default
		Cache:            rdb,
		DB:               db,
	}
	bot.InitBot()

	h := handler.Handler{
		Api:    cfg,
		Logger: zaplog,
		Bot:    &bot,
		Redis:  rdb,
		DB:     db,
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

	////TODO: build host and port from config and call single function
	//err = rh.Run(cfg.Server.Host + ":" + cfg.Server.Port)
	//if err != nil {
	//	cfg.Logger.Fatal("Failed to start server", zap.String("error", err.Error()))
	//}

}

// Migrate the database models
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.City{},
		&models.Area{},
		&models.Product{},
		&models.ProductPrice{},
		&models.AreaProduct{},
		&models.CityProduct{},
	)
}

func populateDatabase(db *gorm.DB) error {
	// Define products with prices
	apple := models.Product{
		Name:        "Apple",
		Description: "Fresh red apples",
		ProductPrices: []models.ProductPrice{
			{Quantity: 1, Price: 2.00},
			{Quantity: 2, Price: 3.00},
			{Quantity: 3, Price: 4.00},
			{Quantity: 4, Price: 5.00},
			{Quantity: 5, Price: 6.00},
		},
	}

	orange := models.Product{
		Name:        "Orange",
		Description: "Juicy oranges",
		ProductPrices: []models.ProductPrice{
			{Quantity: 1, Price: 1.50},
			{Quantity: 2, Price: 2.50},
			{Quantity: 3, Price: 3.50},
			{Quantity: 4, Price: 4.50},
			{Quantity: 5, Price: 5.50},
		},
	}

	banana := models.Product{
		Name:        "Banana",
		Description: "Fresh bananas",
		ProductPrices: []models.ProductPrice{
			{Quantity: 1, Price: 1.20},
			{Quantity: 2, Price: 2.20},
			{Quantity: 3, Price: 3.20},
			{Quantity: 4, Price: 4.20},
			{Quantity: 5, Price: 5.20},
		},
	}

	// Define Jacksonville with its areas and products
	jacksonville := models.City{
		Name: "Jacksonville",
		Products: []models.Product{
			apple, // City-wide available product
			orange,
		},
		Areas: []models.Area{
			{
				Name: "Saint Johns Town Center",
				Products: []models.Product{
					apple,
				},
			},
			{
				Name: "Mandarin",
				Products: []models.Product{
					orange,
				},
			},
			{
				Name: "Orange Park",
				Products: []models.Product{
					orange,
				},
			},
		},
	}

	// Define Chicago with its areas and products
	chicago := models.City{
		Name: "Chicago",
		Products: []models.Product{
			banana,
			apple,
		},
		Areas: []models.Area{
			{
				Name: "Downtown",
				Products: []models.Product{
					banana,
					apple,
				},
			},
			{
				Name: "Lincoln Park",
				Products: []models.Product{
					banana,
				},
			},
		},
	}

	// Define New York with its areas and products
	newYork := models.City{
		Name: "New York",
		Products: []models.Product{
			apple,
			orange,
			banana,
		},
		Areas: []models.Area{
			{
				Name: "Manhattan",
				Products: []models.Product{
					apple,
					orange,
				},
			},
			{
				Name: "Brooklyn",
				Products: []models.Product{
					banana,
					orange,
				},
			},
			{
				Name: "Queens",
				Products: []models.Product{
					apple,
					banana,
				},
			},
		},
	}

	// Add the cities with their areas, products, and product prices to the database
	db.Create(&jacksonville)
	db.Create(&chicago)
	db.Create(&newYork)

	fmt.Println("Database populated successfully with cities, areas, products, and product prices.")
	return nil
}
