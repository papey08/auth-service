package server

import (
	"auth-service/internal/app"
	"auth-service/internal/model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func signIn(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		guid := c.Param("guid")

		tokens, err := a.SignIn(c, model.User{GUID: guid})

		switch {
		case err == nil:
			c.JSON(http.StatusOK, successResponse(tokens))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		}
	}
}

func refresh(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody refreshRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		}

		tokens, err := a.RefreshTokens(c, reqBody.RefreshToken)

		switch {
		case errors.Is(err, model.ExpTokenError):
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		case errors.Is(err, model.NoTokenError):
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		case errors.Is(err, model.DecodeTokenError):
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		case err == nil:
			c.JSON(http.StatusOK, successResponse(tokens))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		}
	}
}
