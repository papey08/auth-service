package app

import (
	"auth-service/internal/model"
	"context"
)

type app struct {
	r *Repo
}

func (a *app) SignIn(ctx context.Context, user model.User) (model.Tokens, error) {
	// TODO: implement
	return model.Tokens{}, nil
}

func (a *app) RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error) {
	// TODO: implement
	return model.Tokens{}, nil
}
