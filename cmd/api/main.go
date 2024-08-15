package main

import (
	"Magaz/internal/config"
	"Magaz/pkg/utils/logger"
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
	cfg.Logger.Info("Starting the application")

	// Start the server

	//bot, err := telego.NewBot("7466095384:AAEg1aQpK6vbp0AWodLGbPALCVhHlKY1_kM", telego.WithDefaultDebugLogger())
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//_ = bot.SetWebhook(&telego.SetWebhookParams{
	//	URL: "https://3e3c-73-192-67-43.ngrok-free.app/webhook",
	//})
	//
	//// Receive information about webhook
	//info, _ := bot.GetWebhookInfo()
	//fmt.Printf("Webhook Info: %+v\n", info)
	//
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//	})
	//})
	//
	//// Create a channel to receive updates
	//updateChannel := make(chan telego.Update, 128) // Buffer size is 128
	//
	////
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
	//r.POST("/webhook", func(c *gin.Context) {
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
	//
	//r.Run(":8080")

}
