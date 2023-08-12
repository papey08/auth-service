package app

import (
	"auth-service/internal/model"
	"context"
	"time"
)

type App interface {
	// SignIn return access and refresh tokens of signed user
	SignIn(ctx context.Context, user model.User) (model.Tokens, error)

	// RefreshTokens return access and refresh tokens if refreshToken is valid
	RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error)
}

func New(r Repo, t Tokenizer) App {
	return &app{
		r: r,
		t: t,
	}
}

type Repo interface {
	// InsertToken adds new refresh token to database
	InsertToken(ctx context.Context, user model.User, token string, expiresAt time.Time) error

	// UpdateToken updates refresh token to a new token
	UpdateToken(ctx context.Context, oldToken, newToken string, expiresAt time.Time) error

	// GetByRefreshToken returns user with given refreshToken and time when token expires
	// or error if refreshToken has expired or not exists
	GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, time.Time, error)

	// RemoveExpiredTokens removes all expired refresh tokens from the database
	RemoveExpiredTokens(ctx context.Context)
}

type Tokenizer interface {
	// NewTokens return new access and refresh tokens
	NewTokens(guid string) (model.Tokens, error)
}
