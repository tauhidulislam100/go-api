package routes

import (
	"gin-api/controllers"
	"gin-api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.Use(middleware.VerifyAuth())
	router.GET("/users", controllers.GetAllUsers())
	router.GET("/users/:id", controllers.GetUser())
	router.PATCH("/users/:id", controllers.UpdateUser())
}
