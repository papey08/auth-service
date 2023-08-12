package app

import (
	repomocks "auth-service/internal/app/repo_mocks"
	tokenmocks "auth-service/internal/app/token_mocks"
	"auth-service/internal/model"
	cryptotools "auth-service/pkg/crypto-tools"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type appTestSuite struct {
	suite.Suite
	repo    *repomocks.Repo
	t       *tokenmocks.Tokenizer
	service App
}

func (s *appTestSuite) SetupSuite() {
	s.repo = new(repomocks.Repo)
	s.t = new(tokenmocks.Tokenizer)
	s.service = New(s.repo, s.t)
}

func (s *appTestSuite) TearDown() {}

type signInMock struct {
	userGUID  string
	newTokens model.Tokens

	user         model.User
	refreshToken string
	expiresAt    time.Time
}

type signInTest struct {
	description string

	givenUser      model.User
	expectedTokens model.Tokens
}

func (s *appTestSuite) TestSignIn() {
	timeNow := time.Now()

	mocks := []signInMock{
		{
			userGUID: "qwerty",
			newTokens: model.Tokens{
				AccessToken: "qwerty-in-jwt",
				RefreshToken: model.RefreshToken{
					Token:     "qwerty-refresh-token",
					ExpiresAt: timeNow.Add(time.Hour),
				},
			},

			user: model.User{
				GUID: "qwerty",
			},
			refreshToken: "qwerty-refresh-token",
			expiresAt:    timeNow.Add(time.Hour),
		},
	}

	tests := []signInTest{
		{
			description: "correct sign in operation",

			givenUser: model.User{
				GUID: "qwerty",
			},
			expectedTokens: model.Tokens{
				AccessToken: "qwerty-in-jwt",
				RefreshToken: model.RefreshToken{
					Token:     cryptotools.StringToBase64("qwerty-refresh-token"),
					ExpiresAt: timeNow.Add(time.Hour),
				},
			},
		},
	}

	for _, m := range mocks {
		s.t.On("NewTokens", m.userGUID).Return(m.newTokens, nil).Once()
		s.repo.On("InsertToken", mock.Anything, m.user, m.refreshToken, m.expiresAt).Return(nil).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		tokens, err := s.service.SignIn(ctx, test.givenUser)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), test.expectedTokens, tokens)
	}
}

type refreshTokensMock struct {
	decodedToken string

	gotUser      model.User
	gotExpiresAt time.Time

	userGUID  string
	newTokens model.Tokens

	newToken     string
	newExpiresAt time.Time
}

type refreshTokenTest struct {
	description string

	givenToken     string
	expectedTokens model.Tokens
	expectedErr    error
}

func (s *appTestSuite) TestRefreshTokens() {
	timeNow := time.Now()

	mocks := []refreshTokensMock{
		{
			decodedToken: "qwerty-refresh-token",

			gotUser: model.User{
				GUID: "qwerty",
			},
			gotExpiresAt: timeNow.Add(time.Hour),

			userGUID: "qwerty",
			newTokens: model.Tokens{
				AccessToken: "new-qwerty-in-jwt",
				RefreshToken: model.RefreshToken{
					Token:     "new-qwerty-refresh-token",
					ExpiresAt: timeNow.Add(time.Hour),
				},
			},

			newExpiresAt: timeNow.Add(time.Hour),
		},
		{
			decodedToken: "qwerty-expired-token",

			gotUser: model.User{
				GUID: "qwerty",
			},
			gotExpiresAt: timeNow,
		},
	}

	tests := []refreshTokenTest{
		{
			description: "correct refresh operation",

			givenToken: cryptotools.StringToBase64("qwerty-refresh-token"),
			expectedTokens: model.Tokens{
				AccessToken: "new-qwerty-in-jwt",
				RefreshToken: model.RefreshToken{
					Token:     cryptotools.StringToBase64("new-qwerty-refresh-token"),
					ExpiresAt: timeNow.Add(time.Hour),
				},
			},
			expectedErr: nil,
		},
		{
			description: "attempt of refresh by expired token",

			givenToken:     cryptotools.StringToBase64("qwerty-expired-token"),
			expectedTokens: model.Tokens{},
			expectedErr:    model.ExpTokenError,
		},
	}

	for _, m := range mocks {
		s.repo.On("GetByRefreshToken", mock.Anything, m.decodedToken).Return(m.gotUser, m.gotExpiresAt, nil).Once()
		s.t.On("NewTokens", m.userGUID).Return(m.newTokens, nil).Once()
		s.repo.On("UpdateToken", mock.Anything, m.decodedToken, mock.Anything, m.newExpiresAt).Return(nil).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		tokens, err := s.service.RefreshTokens(ctx, test.givenToken)
		assert.Equal(s.T(), test.expectedTokens, tokens)
		assert.Equal(s.T(), test.expectedErr, err)
	}
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(appTestSuite))
}
