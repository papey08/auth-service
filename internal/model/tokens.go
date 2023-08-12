package model

import "time"

type RefreshToken struct {
	Token     string
	ExpiresAt time.Time
}

type Tokens struct {
	AccessToken string
	RefreshToken
}
