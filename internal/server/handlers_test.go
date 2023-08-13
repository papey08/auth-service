package server

import (
	"auth-service/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
)

func (s *serverTestSuite) signIn(url string) (tokensData, int, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(s.baseURL+"/auth/v1"+url), nil)
	if err != nil {
		return tokensData{}, 0, err
	}

	var resp response
	code, err := s.getResponse(req, &resp)
	if err != nil {
		return tokensData{}, 0, err
	}

	if resp.Data == nil {
		return tokensData{}, code, nil
	}

	return *resp.Data, code, nil
}

type signInMock struct {
	user   model.User
	tokens model.Tokens
	err    error
}

type signInTest struct {
	description string

	givenURL           string
	expectedTokens     tokensData
	expectedStatusCode int
}

func (s *serverTestSuite) TestSignIn() {
	timeNow := time.Now()

	mocks := []signInMock{
		{
			user: model.User{
				GUID: "qwerty",
			},
			tokens: model.Tokens{
				AccessToken: "qwerty-access-token",
				RefreshToken: model.RefreshToken{
					Token:     "qwerty-refresh-token-in-base64",
					ExpiresAt: timeNow.Add(time.Hour),
				},
			},
			err: nil,
		},
	}

	tests := []signInTest{
		{
			description: "successful sign in",

			givenURL: "/sign-in/qwerty",
			expectedTokens: tokensData{
				AccessToken:  "qwerty-access-token",
				RefreshToken: "qwerty-refresh-token-in-base64",
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, m := range mocks {
		s.app.On("SignIn", mock.Anything, m.user).Return(m.tokens, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.signIn(test.givenURL)
		assert.Equal(s.T(), test.expectedTokens, resp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}
}

func (s *serverTestSuite) refresh(body map[string]any) (tokensData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return tokensData{}, 0, err
	}

	req, err := http.NewRequest(http.MethodPost, s.baseURL+"/auth/v1/refresh", bytes.NewReader(data))
	if err != nil {
		return tokensData{}, 0, err
	}

	req.Header.Add("Content-Type", "application/json")

	var resp response
	code, err := s.getResponse(req, &resp)
	if err != nil {
		return tokensData{}, 0, err
	}

	if resp.Data == nil {
		return tokensData{}, code, nil
	}

	return *resp.Data, code, nil
}

type refreshMock struct {
	refreshToken string
	tokens       model.Tokens
	err          error
}

type refreshTest struct {
	description string

	givenBody          map[string]any
	expectedTokens     tokensData
	expectedStatusCode int
}

func (s *serverTestSuite) TestRefresh() {
	timeNow := time.Now()

	mocks := []refreshMock{
		{
			refreshToken: "qwerty-refresh-token",
			tokens: model.Tokens{
				AccessToken: "qwerty-access-token",
				RefreshToken: model.RefreshToken{
					Token:     "qwerty-refresh-token-in-base64",
					ExpiresAt: timeNow.Add(time.Hour),
				},
			},
			err: nil,
		},
		{
			refreshToken: "qwerty-refresh-token-expired",
			tokens:       model.Tokens{},
			err:          model.ExpTokenError,
		},
	}

	tests := []refreshTest{
		{
			description: "successful refresh operation",

			givenBody: map[string]any{
				"refresh_token": "qwerty-refresh-token",
			},
			expectedTokens: tokensData{
				AccessToken:  "qwerty-access-token",
				RefreshToken: "qwerty-refresh-token-in-base64",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "attempt of refresh by expired token",

			givenBody: map[string]any{
				"refresh_token": "qwerty-refresh-token-expired",
			},
			expectedTokens:     tokensData{},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, m := range mocks {
		s.app.On("RefreshTokens", mock.Anything, m.refreshToken).Return(m.tokens, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.refresh(test.givenBody)
		assert.Equal(s.T(), test.expectedTokens, resp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}
}
