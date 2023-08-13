package server

import "auth-service/internal/model"

// tokenData is a field of response
type tokensData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// response is a struct of any server response. It's fields have pointer type
// for correct marshalling to json
type response struct {
	Data *tokensData `json:"data"`
	Err  *string     `json:"error"`
}

func successResponse(t model.Tokens) response {
	return response{
		Data: &tokensData{
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken.Token,
		},
		Err: nil,
	}
}

func errorResponse(err error) response {
	errStr := err.Error()
	return response{
		Data: nil,
		Err:  &errStr,
	}
}
