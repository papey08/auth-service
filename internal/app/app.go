package app

import (
	"auth-service/internal/model"
	"context"
	"log"
	"time"
)

type app struct {
	r Repo
	t Tokenizer
}

func (a *app) SignIn(ctx context.Context, user model.User) (model.Tokens, error) {
	// generating a new pair of tokens
	tokens, err := a.t.NewTokens(user.GUID)
	if err != nil {
		return model.Tokens{}, model.GenTokenError
	}

	// inserting new tokens to the database
	if err = a.r.InsertToken(ctx, user, tokens.RefreshToken.Token, tokens.RefreshToken.ExpiresAt); err != nil {
		return model.Tokens{}, err
	}
	return tokens, nil
}

func (a *app) RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error) {
	// searching given token in the database
	u, t, err := a.r.GetByRefreshToken(ctx, refreshToken)

	if err != nil {
		log.Println("RefreshTokens: GetByRefreshToken error: ", err.Error())
		return model.Tokens{}, err
	} else if t.Before(time.Now()) {
		return model.Tokens{}, model.ExpTokenError
	}

	// generating a new pair of tokens
	tokens, err := a.t.NewTokens(u.GUID)
	if err != nil {
		return model.Tokens{}, err
	}

	// updating refresh token in the database
	if err = a.r.UpdateToken(ctx, refreshToken, tokens.RefreshToken.Token, tokens.RefreshToken.ExpiresAt); err != nil {
		log.Println("RefreshTokens: UpdateToken error: ", err.Error())
		return model.Tokens{}, err
	}
	return tokens, nil
}
