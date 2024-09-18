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
		admin := api.Group("/admin")
		{
			admin.GET("/", h.AdminHandler())
			admin.GET("/products/getProducts", h.GetProductsAdminHandler())
			admin.POST("/products/add-product", h.PostAdminAddProduct())
			admin.POST("/products/addProdAddr", h.PostAdminAddProductAddr()) //TODO:Rename o address
			admin.GET("/products/getProdAddr", h.AdminGetProductAddr())      //TODO:Rename o address
		}

		empl := api.Group("/empl")
		{
			//empl.GET("/orders", h.HEmployeeHandler())
			//empl.GET("/orders", h.EmployeeHandler())

			//TODO: Move later to api section , since admin might use for communication in future too
			empl.GET("/ws", h.Upgrade())

			empl.GET("/orders", h.EmployeeHandlerTest())
			empl.GET("/orders/address", h.GetOrderAddressHandler())
			empl.POST("/orders/address/assign", h.PostOrderAddressHandler())
			empl.POST("/orders/release/:orderId", h.ReleaseOrderHandler())

		}

		//General GET calls
		api.GET("/get/images/:image", h.ServeImage())

		//POST calls
		api.POST("/bot/telegram", h.BotRequestHandler())
	}

	// Serve static files from the "static" directory
	//TODO: pass the path from config
	router.Static("frontend/v1/static", "./frontend/v1/static")
	router.Static("frontend/v4/static", "./frontend/v4/static")

	return router
}
