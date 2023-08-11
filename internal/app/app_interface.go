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

func NewApp(r *Repo) App {
	return &app{r: r}
}

type Repo interface {
	// SetSession updates refresh token in database
	SetSession(ctx context.Context, user model.User, refreshToken string, expiresAt time.Time) error

	// GetByRefreshToken returns user with given refreshToken
	GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error)
}
