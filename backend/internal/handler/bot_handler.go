package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	"net/http"
)

// BotRequestHandler handles the incoming requests from the Telegram bot
func (h *Handler) BotRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var update telego.Update

		// Bind the JSON to the update struct
		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update structure"})
			return
		}
		//TODO: Do some other logic here , validating , or possibly connect with other bot to tech support

		// Send the update to the channel for processing
		h.Bot.UpdatesChan <- update

		// Respond with status OK
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
