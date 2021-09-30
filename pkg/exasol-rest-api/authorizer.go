package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"net/http"
)

// Authorizer is responsible for the API users' authorization
type Authorizer interface {
	// Authorize a token
	Authorize(request *http.Request) error
}

// TokenAuthorizer is a token-based implementation of the Authorizer
type TokenAuthorizer struct {
	AllowedTokens map[string]bool
}

func (auth *TokenAuthorizer) Authorize(request *http.Request) error {
	tokens := request.Header["Authorization"]

	authorized := false
	for _, token := range tokens {
		if len(token) < 30 {
			errorLogger.Printf("attempt to access API with a token of a wrong length")
			return fmt.Errorf(error_reporting_go.ExaError("E-ERA-23").
				Message("an authorization token has invalid length: {{length|uq}}.").
				Parameter("length", len(token)).
				Mitigation("please only use tokens with the length longer or equal to 30.").Error())
		}

		if auth.AllowedTokens[token] {
			authorized = true
			break
		}
	}

	if !authorized {
		errorLogger.Printf("attempt to access API with an invalid token")
		return fmt.Errorf(error_reporting_go.ExaError("E-ERA-22").
			Message("an authorization token is missing or wrong.").
			Mitigation("please make sure you provided a valid token.").Error())
	} else {
		return nil
	}
}
