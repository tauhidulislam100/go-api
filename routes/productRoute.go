package routes

import (
	"gin-api/controllers"
	"gin-api/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRouter(router *gin.Engine) {
	router.Use(middleware.VerifyAuth())
	router.POST("products", controllers.CreateProduct())
	router.GET("products", controllers.GetAllProducts())
	router.GET("products/:id", controllers.GetProduct())
	router.PATCH("products/:id", controllers.UpdateProduct())
}
