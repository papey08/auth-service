package tokenizer

import (
	"auth-service/internal/model"
	"github.com/Pallinder/go-randomdata"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	accessTokenExpiresAt  = time.Minute * 15
	refreshTokenExpiresAt = time.Hour * 24 * 30

	refreshTokenLength = 64
)

type Tokenizer struct {
	signKey string
}

func New(s string) *Tokenizer {
	return &Tokenizer{
		signKey: s,
	}
}

// NewAccessToken return access token with encrypted guid
func (t *Tokenizer) newAccessToken(data string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = data
	claims["exp"] = time.Now().Add(accessTokenExpiresAt).Unix()
	return token.SignedString([]byte(t.signKey))
}

// NewRefreshToken generates new refresh token
func (t *Tokenizer) newRefreshToken() model.RefreshToken {
	return model.RefreshToken{
		Token:     randomdata.Alphanumeric(refreshTokenLength),
		ExpiresAt: time.Now().Add(refreshTokenExpiresAt),
	}
}

func (t *Tokenizer) NewTokens(guid string) (model.Tokens, error) {
	var tokens model.Tokens
	var err error

	if tokens.AccessToken, err = t.newAccessToken(guid); err != nil {
		return model.Tokens{}, model.GenTokenError
	}

	tokens.RefreshToken = t.newRefreshToken()
	return tokens, nil
}
