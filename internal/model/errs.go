package model

import "errors"

var (
	GenTokenError    = errors.New("unable to generate a pair of tokens")
	RepoError        = errors.New("something wrong with repository")
	ExpTokenError    = errors.New("token has expired")
	NoTokenError     = errors.New("no required token in database")
	TokenCryptError  = errors.New("unable to encrypt refresh token")
	DecodeTokenError = errors.New("unable to decode refresh token")
)
