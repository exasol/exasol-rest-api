package exasol_rest_api_test

import (
	"github.com/stretchr/testify/suite"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"net/http"
	"testing"
)

type AuthorizerSuite struct {
	suite.Suite
}

func TestAuthorizerSuite(t *testing.T) {
	suite.Run(t, new(AuthorizerSuite))
}

func (suite *AuthorizerSuite) TestSingleTokenAuthorized() {
	authorizer := exasol_rest_api.TokenAuthorizer{
		AllowedTokens: map[string]bool{
			"nqtfbaD34DSHhzUHIN2VYoTCo6ULne": true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	onError(err)
	req.Header.Set("Authorization", "nqtfbaD34DSHhzUHIN2VYoTCo6ULne")
	suite.Equal(nil, authorizer.Authorize(req))
}

func (suite *AuthorizerSuite) TestMultipleTokensAuthorized() {
	authorizer := exasol_rest_api.TokenAuthorizer{
		AllowedTokens: map[string]bool{
			"nqtfbaD34DSHhzUHIN2VYoTCo6ULne":    true,
			"OO2SJQ8CSqSKjvU8DPqWuL0OZ2ewwi45":  true,
			"yGrdSdoFqq3plkrsZuKSpF6F7f5s8q234": true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	onError(err)
	req.Header.Set("Authorization", "FN4xjKW40LlgL7DonFnHwozUV3YAeO")
	req.Header.Set("Authorization", "OO2SJQ8CSqSKjvU8DPqWuL0OZ2ewwi45")
	suite.Equal(nil, authorizer.Authorize(req))
}

func (suite *AuthorizerSuite) TestNotAuthorized() {
	authorizer := exasol_rest_api.TokenAuthorizer{
		AllowedTokens: map[string]bool{
			"nqtfbaD34DSHhzUHIN2VYoTCo6ULne": true,
			"OO2SJQ8CSqSKjvU8DPqWuL0OZ2ewwi": true,
			"yGrdSdoFqq3plkrsZuKSpF6F7f5s8q": true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	onError(err)
	req.Header.Set("Authorization", "3J90XAv9loMIXzQdfYmtJrHAbopPsc")
	req.Header.Set("Authorization", "9O0Ynm21G5tAZckiXrQqOTA0br63FW")
	suite.EqualError(authorizer.Authorize(req),
		"E-ERA-22: an authorization token is missing or wrong. please make sure you provided a valid token.")
}

func (suite *AuthorizerSuite) TestInvalidTokenLength() {
	authorizer := exasol_rest_api.TokenAuthorizer{
		AllowedTokens: map[string]bool{
			"abc": true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	onError(err)
	req.Header.Set("Authorization", "abc")
	suite.EqualError(authorizer.Authorize(req),
		"E-ERA-23: an authorization token has invalid length: 3. please only use tokens with the length longer or equal to 30.")
}
