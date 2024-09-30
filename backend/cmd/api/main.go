package main

import (
	"Magaz/backend/internal/config"
	"Magaz/backend/internal/handler"
	"Magaz/backend/internal/router"
	"Magaz/backend/internal/storage/models"
	"Magaz/backend/internal/system/websocket"
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

	cfg, err := config.LoadAPIConfig()
	if err != nil {
		log.Fatalf("Failed to load API configs: %v", err)
	}

	h := handler.NewHandler(cfg)

	zaplog, _ := logger.InitLogger(cfg.Env)
	h.Logger = zaplog

	rdb, rdberr := redis.InitRedisClient(&cfg.Redis) //TODO: assign to var
	if rdberr != nil {
		zaplog.Fatal("Failed to connect to Redis", zap.String("error", rdberr.Error()))
	}
	h.Redis = rdb

	db, dberr := postgres.Connect(cfg.Database)
	if dberr != nil {
		zaplog.Fatal("Failed to connect to DB", zap.String("error", dberr.Error()))

	}
	h.DB = db

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
	h.Session = store

	//hub := sse.NewSSEHub(db, zaplog)
	//go hub.Run()
	wscon := ws.NewManager(zaplog, db)
	h.WS = wscon

	//TODO: passing to handler initialized clients like redis and db . Pass handler instead ?
	bot := telegram.Bot{
		Config:           &cfg.Bot,
		Logger:           zaplog,
		UpdateChanBuffer: 128, // Buffer size is 128 default
		Cache:            rdb,
		DB:               db,
		WS:               wscon,
		//Hub:              hub,
	}
	bot.InitBot()
	h.Bot = &bot

	rh := router.SetupRouter(h)

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
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
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
		&models.DeclinedOrder{},
		&models.Card{},
		&models.Crypto{},
	)
}

func populateTestData(db *gorm.DB) {
	// Define some test card data
	cards := []models.Card{
		{
			BankName:   "Chase Bank",
			BankURL:    "https://www.chase.com",
			UserID:     "john_doe",
			Password:   "password123",
			CardNumber: "4111111111111111",
			QuickPay:   "enabled",
			FirstName:  "John",
			LastName:   "Doe",
			ExpireDate: "12/25",
			CVV:        "123",
			CardType:   "Visa",
			Balance:    10000, // $100.00 balance
			Active:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			BankName:   "Bank of America",
			BankURL:    "https://www.bankofamerica.com",
			UserID:     "jane_doe",
			Password:   "password456",
			CardNumber: "4222222222222",
			QuickPay:   "enabled",
			FirstName:  "Jane",
			LastName:   "Doe",
			ExpireDate: "11/24",
			CVV:        "456",
			CardType:   "MasterCard",
			Balance:    15000, // $150.00 balance
			Active:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			BankName:   "Wells Fargo",
			BankURL:    "https://www.wellsfargo.com",
			UserID:     "bob_smith",
			Password:   "password789",
			CardNumber: "4333333333333",
			QuickPay:   "disabled",
			FirstName:  "Bob",
			LastName:   "Smith",
			ExpireDate: "10/23",
			CVV:        "789",
			CardType:   "Discover",
			Balance:    5000, // $50.00 balance
			Active:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	// Insert the test cards into the database
	for _, card := range cards {
		if err := db.Create(&card).Error; err != nil {
			log.Printf("Failed to insert card: %v", err)
		} else {
			log.Printf("Inserted card: %s %s", card.FirstName, card.LastName)
		}
	}
}
