package tokenizer

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccessToken(t *testing.T) {
	tests := []string{
		"abcdef",
		"абвгде",
		"qwerty1234_0987#$%",
	}
	tt := New("secret-token-key")

	for _, test := range tests {
		tokenStr, err := tt.newAccessToken(test)
		assert.NoError(t, err)

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (any, error) {
			return []byte("secret-token-key"), nil
		})
		assert.NoError(t, err)

		claims, ok := token.Claims.(*jwt.StandardClaims)
		assert.True(t, ok)

		assert.Equal(t, test, claims.Subject)
	}
}
