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
		case errors.Is(err, model.RepoError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		case errors.Is(err, model.GenTokenError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		case errors.Is(err, model.TokenCollisionError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		case err == nil:
			c.JSON(http.StatusOK, successResponse(tokens))
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
		case errors.Is(err, model.RepoError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		case errors.Is(err, model.GenTokenError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		case err == nil:
			c.JSON(http.StatusOK, successResponse(tokens))
		}
	}
}
