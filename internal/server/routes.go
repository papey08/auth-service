package server

import "github.com/gin-gonic/gin"

func routes(r *gin.RouterGroup) {
	r.GET("/sign-in/:guid") // TODO: make handler for the first route
	r.POST("/refresh")      // TODO: make handler for the second route
}
