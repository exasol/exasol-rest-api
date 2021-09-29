package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"net/http"
)

// Authorizer is responsible for the API users' authorization
type Authorizer interface {
	authorize(request *http.Request) error
}

// TokenAuthorizer is a token-based implementation of the Authorizer
type TokenAuthorizer struct {
	AllowedToken string
}

func (auth *TokenAuthorizer) authorize(request *http.Request) error {
	tokens := request.Header["Authorization"]

	authorized := false
	for _, token := range tokens {
		if token == auth.AllowedToken {
			authorized = true
		}
	}

	if !authorized {
		return fmt.Errorf(error_reporting_go.ExaError("E-ERA-22").
			Message("an authorization token is missing or wrong.").
			Mitigation("please make sure you provided a valid token.").Error())
	} else {
		return nil
	}
}
