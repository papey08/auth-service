package app

import (
	"auth-service/internal/model"
	"context"
)

type App interface {
	// SignIn return access and refresh tokens of signed user
	SignIn(ctx context.Context, user model.User) (model.Tokens, error)

	// RefreshTokens return access and refresh tokens if refreshToken is valid
	RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error)
}

func NewApp(r Repo, t Tokenizer) App {
	return &app{
		r: r,
		t: t,
	}
}

type Repo interface {
	// InsertToken adds new refresh token to database
	InsertToken(ctx context.Context, user model.User, token string) error

	// UpdateToken updates refresh token to a new token
	UpdateToken(ctx context.Context, oldToken, newToken string) error

	// GetByRefreshToken returns user with given refreshToken
	GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error)
}

type Tokenizer interface {
	// NewTokens return new access and refresh tokens
	NewTokens(guid string) (model.Tokens, error)

	// TokenValidate checks if token has not expired and it is valid
	TokenValidate(tokenStr string) error
}
