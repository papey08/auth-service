package tokenizer

import (
	"auth-service/internal/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	accessTokenExpiresAt  = time.Minute * 15
	refreshTokenExpiresAt = time.Hour * 24 * 30
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
func (t *Tokenizer) newAccessToken(guid string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTokenExpiresAt).Unix(),
		Subject:   guid,
	}).SignedString([]byte(t.signKey))
}

// NewRefreshToken returns refresh token
func (t *Tokenizer) newRefreshToken() (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(refreshTokenExpiresAt).Unix(),
	}).SignedString([]byte(t.signKey))
}

func (t *Tokenizer) NewTokens(guid string) (model.Tokens, error) {
	var tokens model.Tokens
	var err error
	if tokens.AccessToken, err = t.newAccessToken(guid); err != nil {
		return model.Tokens{}, model.GenTokenError
	}
	if tokens.RefreshToken, err = t.newRefreshToken(); err != nil {
		return model.Tokens{}, model.GenTokenError
	}
	return tokens, nil
}

func TokenValidate(tokenStr string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return model.InvalidTokenError
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().Before(expirationTime) {
			return nil
		} else {
			return model.ExpTokenError
		}
	} else {
		return model.InvalidTokenError
	}
}
