package model

import "errors"

var (
	GenTokenError       = errors.New("unable to generate a pair of tokens")
	RepoError           = errors.New("something wrong with repository")
	ExpTokenError       = errors.New("token has expired")
	NoTokenError        = errors.New("no required token in database")
	TokenCollisionError = errors.New("required token is already in database")
)
