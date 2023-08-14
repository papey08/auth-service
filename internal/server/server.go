package server

import (
	"auth-service/internal/app"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewHTTPServer creates http.Server with routes
func NewHTTPServer(a app.App, host string, port int) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("auth/v1")

	routes(api, a)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
}
