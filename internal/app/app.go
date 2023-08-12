package app

import (
	"auth-service/internal/model"
	cryptotools "auth-service/pkg/crypto-tools"
	"context"
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

	// changing refresh token format to base64
	tokens.RefreshToken.Token = cryptotools.StringToBase64(tokens.RefreshToken.Token)
	return tokens, nil
}

func (a *app) RefreshTokens(ctx context.Context, refreshToken string) (model.Tokens, error) {

	// decoding given refresh token from base64
	decodedRefreshToken, err := cryptotools.Base64ToString(refreshToken)
	if err != nil {
		return model.Tokens{}, model.DecodeTokenError
	}

	// searching given token in the database
	u, t, err := a.r.GetByRefreshToken(ctx, decodedRefreshToken)

	if err != nil {
		return model.Tokens{}, err
	} else if t.Before(time.Now()) {
		return model.Tokens{}, model.ExpTokenError
	}

	// generating a new pair of tokens
	tokens, err := a.t.NewTokens(u.GUID)
	if err != nil {
		return model.Tokens{}, err
	}

	// encrypting new refresh token
	encryptedNewRefreshToken, err := cryptotools.GenerateBcryptHash(tokens.RefreshToken.Token)
	if err != nil {
		return model.Tokens{}, model.TokenCryptError
	}

	// updating refresh token in the database
	if err = a.r.UpdateToken(ctx,
		decodedRefreshToken,
		encryptedNewRefreshToken,
		tokens.RefreshToken.ExpiresAt); err != nil {
		return model.Tokens{}, err
	}
	return tokens, nil
}
