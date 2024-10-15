package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"
	"golang.ngrok.com/ngrok"
	cfg "golang.ngrok.com/ngrok/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tg/internal/bot"
	"tg/internal/config"
	"tg/internal/utlis/logger"
	"tg/pkg/utils/service"
	"time"
)

type AppConfig struct {
	BotSrv *bot.BotService
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := AppConfig{}

	zlog, err := logger.InitLogger("dev") //TODO: get propper env definition
	if err != nil {
		log.Print("failed to initialize logger: %w", err)
	}

	botCfg, err := config.LoadBotConfig()
	if err != nil {
		zlog.Fatal("Failed to load bot config")
	}

	///TODO: Move this to NewProducer init call

	// Start ngrok tunnel
	tun, err := ngrok.Listen(context.Background(),
		cfg.HTTPEndpoint(cfg.WithForwardsTo(":8080")),
		ngrok.WithAuthtoken("46MTSNmoEY38FEtUhXiqH_2ZwvfvSfDxzSKFuvpdJyk"), //TODO: retrieve from cfg
	)
	if err != nil {
		zlog.Fatal("Failed to start server")

	}
	zlog.Info("Ngrok tunnel established", zap.String("url", tun.URL()))

	//TODO: refactor
	botSrv := bot.NewBotService(botCfg)
	botSrv.URL = tun.URL()
	app.BotSrv = botSrv

	serviceMng := service.NewServiceManager()
	serviceMng.RegisterService("bot-service", botSrv)
	err = serviceMng.EnableService("bot-service")
	if err != nil {
		zlog.Error("Failed to enable bot service", zap.Error(err))
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Serve(tun); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zlog.Fatal("Failed to listen and serve", zap.String("error", err.Error()))

		}
	}()

	<-quit
	zlog.Info("Shutting down bot server...")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		zlog.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	// Stop the bot service
	err = serviceMng.DisableService("bot-service")
	if err != nil {
		zlog.Fatal("Failed to stop bot service", zap.Error(err))
	}

}

func (cfg *AppConfig) routes() *gin.Engine {
	router := gin.Default()

	for _, tokenConfig := range cfg.BotSrv.Config.Bot.Tokens {
		whPath := cfg.BotSrv.Config.Bot.WebhookBasePath + "/" + tokenConfig.Type

		// Capture tokenConfig within the closure
		token := tokenConfig.Token

		router.POST(whPath, func(c *gin.Context) {
			queryToken := c.Query("token") // Token from the query params
			if queryToken != token {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
				return
			}

			var update telego.Update
			if err := c.ShouldBindJSON(&update); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update structure"})
				return
			}

			if updateChan, ok := cfg.BotSrv.Updates[token]; ok {
				updateChan <- update
				c.Status(http.StatusOK)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
			}
		})
	}

	return router
}
