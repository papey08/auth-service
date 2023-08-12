package crypto_tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	tests := []string{
		"abcdef",
		"абвгде",
		"qwerty1234_0987#$%",
	}

	for _, orig := range tests {
		encodedString := StringToBase64(orig)
		assert.NotEqual(t, orig, encodedString)

		decodedString, err := Base64ToString(encodedString)
		assert.Equal(t, orig, decodedString)
		assert.NoError(t, err)
	}
}

func TestBcryptHash(t *testing.T) {
	tests := []string{
		"abcdef",
		"абвгде",
		"qwerty1234_0987#$%",
	}

	for _, orig := range tests {
		hash, err := GenerateBcryptHash(orig)
		assert.NotEqual(t, orig, hash)
		assert.NoError(t, err)

		res := CheckHash(hash, orig)
		assert.True(t, res)
	}
}
