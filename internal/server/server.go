package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewHTTPServer creates http.Server with routes
func NewHTTPServer(host string, port int) *http.Server {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	api := router.Group("auth/v1")
	routes(api)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
}
