package routes

import (
	"chatapp/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(rg *gin.RouterGroup) {
	rg.POST("/signup", auth.Signup())
	rg.POST("/login", auth.Login())
	rg.GET("/user/:id", auth.Authenticate(), auth.GetUser())
	
}
