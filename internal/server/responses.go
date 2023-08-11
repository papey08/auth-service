package server

import (
	"auth-service/internal/model"
	"github.com/gin-gonic/gin"
)

type response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func successResponse(tokens model.Tokens) *gin.H {
	return &gin.H{
		"data": response{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		"error": nil,
	}
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
