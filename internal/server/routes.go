package server

import (
	"auth-service/internal/app"
	"github.com/gin-gonic/gin"
)

func routes(r *gin.RouterGroup, a app.App) {
	r.GET("/sign-in/:guid", signIn(a))
	r.POST("/refresh", refresh(a))
}
