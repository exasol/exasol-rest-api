package exasol_rest_api

import (
	"os"
	"strconv"
	"strings"

	exaerror "github.com/exasol/error-reporting-go"
)

func (applicationProperties *ApplicationProperties) setValuesFromEnvironmentVariables() {
	applicationProperties.setExasolUser()
	applicationProperties.setExasolPassword()
	applicationProperties.setExasolHost()
	applicationProperties.setServerAddress()
	applicationProperties.setExasolPort()
	applicationProperties.setAPITokens()
}

func (applicationProperties *ApplicationProperties) setExasolUser() {
	exasolUser := os.Getenv(ExasolUserKey)
	if exasolUser != "" {
		applicationProperties.ExasolUser = exasolUser
	}
}

func (applicationProperties *ApplicationProperties) setExasolPassword() {
	exasolPassword := os.Getenv(ExasolPasswordKey)
	if exasolPassword != "" {
		applicationProperties.ExasolPassword = exasolPassword
	}
}

func (applicationProperties *ApplicationProperties) setExasolHost() {
	exasolHost := os.Getenv(ExasolHostKey)
	if exasolHost != "" {
		applicationProperties.ExasolHost = exasolHost
	}
}

func (applicationProperties *ApplicationProperties) setServerAddress() {
	serverAddress := os.Getenv(ApplicationServerKey)
	if serverAddress != "" {
		applicationProperties.ApplicationServer = serverAddress
	}
}

func (applicationProperties *ApplicationProperties) setExasolPort() {
	exasolPort := os.Getenv(ExasolPortKey)
	if exasolPort != "" {
		port, err := strconv.Atoi(exasolPort)
		if err != nil {
			logEnvironmentVariableParsingError(ExasolPortKey, err)
		} else {
			applicationProperties.ExasolPort = port
		}
	}
}

func logEnvironmentVariableParsingError(variableName string, err error) {
	errorLogger.Print(exaerror.New("E-ERA-5").
		Message("cannot parse environment variable "+variableName+". {{error|uq}}").
		Parameter("error", err.Error()).String())
}

func (applicationProperties *ApplicationProperties) setAPITokens() {
	apiTokens := os.Getenv(APITokensKey)
	if apiTokens != "" {
		applicationProperties.APITokens = strings.Split(apiTokens, ",")
	}
}
