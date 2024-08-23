package router

import (
	"Magaz/internal/handler"
	"Magaz/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	//router.Use(middleware.LogDetailedRequestsMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.SessionMiddleware(h.Session))

	//router.POST(h.Api.Bot.WebhookPath, h.BotRequestHandler())
	api := router.Group("/api")
	{

		api.GET("/login", h.GETLoginHandler())
		api.POST("/login", h.POSTLoginHandler())

		api.Use(middleware.AuthRequired(h.Session))
		{
			api.GET("/admin", h.AdminHandler())
		}

		//api.GET("/employee", h.EmployeeHandler())

		// Define your routes here
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		//POST calls
		api.POST("/bot/telegram", h.BotRequestHandler())
	}

	return router
}
