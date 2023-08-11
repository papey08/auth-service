package app

import (
	"auth-service/internal/model"
	"context"
)

type app struct {
	r Repo
	t Tokenizer
}

func (a *app) SignIn(ctx context.Context, user model.User) (model.Tokens, error) {
	tokens, err := a.t.NewTokens(user.GUID)
	if err != nil {
		return model.Tokens{}, model.GenTokenError
	}
	if err = a.r.InsertToken(ctx, user, tokens.RefreshToken); err != nil {
		return model.Tokens{}, model.RepoError
	}
	return tokens, nil
}

func (a *app) RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error) {
	if err := a.t.TokenValidate(refreshToken); err != nil {
		return model.Tokens{}, err
	}

	u, err := a.r.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return model.Tokens{}, err
	}

	tokens, err := a.t.NewTokens(u.GUID)
	if err != nil {
		return model.Tokens{}, err
	}

	if err = a.r.UpdateToken(ctx, refreshToken, tokens.RefreshToken); err != nil {
		return model.Tokens{}, err
	}
	return tokens, nil
}
