package model

import "errors"

var (
	GenTokenError     = errors.New("unable to generate a pair of tokens")
	RepoError         = errors.New("something wrong with repository")
	InvalidTokenError = errors.New("token is invalid")
	ExpTokenError     = errors.New("token has expired")
	TokenParseError   = errors.New("unable to parse token")
	NoTokenError      = errors.New("no required token in database")
)
