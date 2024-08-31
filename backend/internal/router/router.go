package router

import (
	"Magaz/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	//router.LoadHTMLGlob("web/static/test/*")

	//router.Use(middleware.LogDetailedRequestsMiddleware())
	router.Use(gin.Recovery())
	//router.Use(middleware.SessionMiddleware(h.Session))

	//router.POST(h.Api.Bot.WebhookPath, h.BotRequestHandler())
	api := router.Group("/api")
	{

		//TODO: set "/" roting to login page
		api.GET("/login", h.GETLoginHandler())
		api.POST("/login", h.POSTLoginHandler())

		//TODO:return back after testing complete
		//api.Use(middleware.AuthRequired(h.Session))
		//{
		//	api.GET("/admin", h.AdminHandler())
		//}
		api.GET("/admin", h.AdminHandler())
		api.POST("/admin/products/add-product", h.PostAdminAddProduct())
		api.POST("/admin/products/addItem", h.PostAdminAddProductItem())
		api.GET("/admin/products/getItem", h.AdminGetProductItem())
		//api.GET("/admintest", h.AdminHandlerTest()) //TODO: remove after testing

		api.GET("/employee", h.EmployeeHandler())

		// Define your routes here
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		//POST calls
		api.POST("/bot/telegram", h.BotRequestHandler())
	}

	// Serve static files from the "static" directory
	//TODO: pass the path from config
	router.Static("frontend/static", "./frontend/static")

	return router
}
