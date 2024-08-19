package exasol_rest_api

import (
	"net/http"

	exaerror "github.com/exasol/error-reporting-go"
)

// Authorizer is responsible for the API users' authorization.
type Authorizer interface {
	// Authorize a token
	Authorize(request *http.Request) error
}

// TokenAuthorizer is a token-based implementation of the Authorizer.
type TokenAuthorizer struct {
	AllowedTokens map[string]bool
}

// Authorize validates a user request.
// [impl->dsn~execute-query-headers~1]
// [impl->dsn~get-tables-headers~1]
func (auth *TokenAuthorizer) Authorize(request *http.Request) error {
	tokens := request.Header["Authorization"]

	authorized := false
	for _, token := range tokens {
		if len(token) < 30 {
			errorLogger.Print("attempt to access API with a token of a wrong length")
			return exaerror.New("E-ERA-23").
				Message("an authorization token has invalid length: {{length|uq}}.").
				Parameter("length", len(token)).
				Mitigation("please only use tokens with the length longer or equal to 30.")
		}

		if auth.AllowedTokens[token] {
			authorized = true
			break
		}
	}

	if !authorized {
		errorLogger.Print("attempt to access API with an invalid token")
		return exaerror.New("E-ERA-22").
			Message("an authorization token is missing or wrong.").
			Mitigation("please make sure you provided a valid token.")
	} else {
		return nil
	}
}
