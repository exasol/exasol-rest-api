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
			"abc": true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	onError(err)
	req.Header.Set("Authorization", "abc")
	suite.Equal(nil, authorizer.Authorize(req))
}

func (suite *AuthorizerSuite) TestMultipleTokensAuthorized() {
	authorizer := exasol_rest_api.TokenAuthorizer{
		AllowedTokens: map[string]bool{
			"abc": true,
			"bca": true,
			"acb": true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	onError(err)
	req.Header.Set("Authorization", "kkk")
	req.Header.Set("Authorization", "bca")
	suite.Equal(nil, authorizer.Authorize(req))
}

func (suite *AuthorizerSuite) TestNotAuthorized() {
	authorizer := exasol_rest_api.TokenAuthorizer{
		AllowedTokens: map[string]bool{
			"abc": true,
			"bca": true,
			"acb": true,
		},
	}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	onError(err)
	req.Header.Set("Authorization", "kkk")
	req.Header.Set("Authorization", "bbb")
	suite.EqualError(authorizer.Authorize(req),
		"E-ERA-22: an authorization token is missing or wrong. please make sure you provided a valid token.")
}
