package server

import (
	"auth-service/internal/app"
	"github.com/gin-gonic/gin"

	_ "auth-service/docs/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func routes(r *gin.RouterGroup, a app.App) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/sign-in/:guid", signIn(a))
	r.POST("/refresh", refresh(a))
}
