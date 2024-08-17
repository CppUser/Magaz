package router

import (
	"Magaz/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.POST(h.Api.Bot.WebhookPath, h.BotRequestHandler())

	// Define your routes here
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//updateChannel := make(chan telego.Update, 128) // Buffer size is 128
	//
	////// Start a goroutine to handle updates
	//go func() {
	//	for update := range updateChannel {
	//		// Here you can handle the update
	//		if update.Message != nil {
	//			log.Printf("Received message: %s", update.Message.Text)
	//			// Add your update handling logic here
	//		}
	//	}
	//}()
	//
	//router.POST(cfg.Bot.WebhookPath, func(c *gin.Context) {
	//	var update telego.Update
	//	if err := c.BindJSON(&update); err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//		return
	//	}
	//
	//	// Send the update to the channel for processing
	//	updateChannel <- update
	//
	//	// Respond with status OK
	//	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	//})

	// Additional routes can be added here
	// router.GET("/some-route", someHandlerFunction)

	return router
}
