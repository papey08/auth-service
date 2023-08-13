package server

import (
	"auth-service/internal/app"
	"auth-service/internal/model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary		Аутентификация по guid
// @Description	Возвращает пару из Access и Refresh токенов
// @Produce		json
// @Param			guid	path		string			true	"GUID пользователя"
// @Success		200		{object}	server.response	"Успешная аутентификация"
// @Failure		500		{object}	server.response	"Проблемы на стороне сервера"
// @Router			/sign-in/{guid} [post]
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

// @Summary		Обновление токенов
// @Description	Возвращает новые Access и Refresh токены
// @Access			json
// @Produce		json
// @Param			input	body		server.refreshRequest	true	"Refresh токен в JSON"
// @Success		200		{object}	server.response			"Успешное обновление токенов"
// @Failure		401		{object}	server.response			"Срок действия токена истёк"
// @Failure		500		{object}	server.response			"Проблемы на стороне сервера"
// @Router			/refresh [post]
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
