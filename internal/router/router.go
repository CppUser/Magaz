package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Define your routes here
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Additional routes can be added here
	// router.GET("/some-route", someHandlerFunction)

	return router
}
