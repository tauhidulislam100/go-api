package routes

import (
	"gin-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
	router.POST("/auth/register", controllers.CreateUser())
	router.POST("/auth/login", controllers.LoginUser())
}
