package main

import (
	"gin-api/configs"
	"gin-api/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	configs.ConnectDB()
}

func main() {
	router := gin.Default()
	routes.UserRoute(router)
	routes.AuthRoute(router)
	routes.ProductRouter(router)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server is running",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
