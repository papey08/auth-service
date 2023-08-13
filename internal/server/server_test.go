package server

import (
	mocks "auth-service/internal/server/app_mocks"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serverTestSuite struct {
	suite.Suite
	app     *mocks.App
	client  *http.Client
	srv     *http.Server
	baseURL string
}

func httpServerTestSuiteInit(s *serverTestSuite) {
	s.app = new(mocks.App)
	s.srv = NewHTTPServer(s.app, "localhost", 8080)
	testServer := httptest.NewServer(s.srv.Handler)
	s.client = testServer.Client()
	s.baseURL = testServer.URL
}

func (s *serverTestSuite) SetupSuite() {
	httpServerTestSuiteInit(s)
}

func (s *serverTestSuite) TearDownSuite() {
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failde to shutdown server: %v", err)
	}
}

func (s *serverTestSuite) getResponse(req *http.Request, out any) (int, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("unexpected error: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("unable to read response: %w", err)
	}
	_ = json.Unmarshal(respBody, out)
	return resp.StatusCode, nil
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}
