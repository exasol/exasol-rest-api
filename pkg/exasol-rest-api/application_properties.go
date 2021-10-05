package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"os"
)

const APITokensKey string = "API_TOKENS"
const ApplicationServerKey string = "SERVER_ADDRESS"
const ExasolUserKey string = "EXASOL_USER"
const ExasolPasswordKey string = "EXASOL_PASSWORD"
const ExasolHostKey string = "EXASOL_HOST"
const ExasolPortKey string = "EXASOL_PORT"
const ExasolWebsocketAPIVersionKey string = "EXASOL_WEBSOCKET_API_VERSION"
const EncryptionKey string = "EXASOL_ENCRYPTION"
const UseTLSKey string = "EXASOL_TLS"

// ApplicationProperties for Exasol REST API service.
type ApplicationProperties struct {
	APITokens                 []string `yaml:"api-tokens"`
	ApplicationServer         string   `yaml:"server-address"`
	ExasolUser                string   `yaml:"exasol-user"`
	ExasolPassword            string   `yaml:"exasol-password"`
	ExasolHost                string   `yaml:"exasol-host"`
	ExasolPort                int      `yaml:"exasol-port"`
	ExasolWebsocketAPIVersion int      `yaml:"exasol-websocket-api-version"`
	Encryption                bool     `yaml:"encryption"`
	UseTLS                    bool     `yaml:"use-tls"`
}

// GetApplicationProperties creates an application properties.
func GetApplicationProperties() *ApplicationProperties {
	properties := readApplicationProperties()
	err := properties.validate()
	if err != nil {
		panic(error_reporting_go.ExaError("E-ERA-7").Message("application properties validation failed. {{error|uq}}").
			Parameter("error", err.Error()).Error())
	}
	return &properties
}

func readApplicationProperties() ApplicationProperties {
	properties := readApplicationPropertiesFromFile()
	properties.setValuesFromEnvironmentVariables()
	properties.fillMissingWithDefaultValues()
	return properties
}

func readApplicationPropertiesFromFile() ApplicationProperties {
	propertiesFilePath := os.Getenv("APPLICATION_PROPERTIES_PATH")
	properties, err := getPropertiesFromFile(propertiesFilePath)
	if err != nil {
		errorLogger.Printf(error_reporting_go.ExaError("E-ERA-6").
			Message("cannot read properties from specified file: {{file path}}. {{error|uq}}").
			Parameter("file path", propertiesFilePath).Parameter("error", err.Error()).String())
		return ApplicationProperties{}
	} else {
		return properties
	}
}

func (applicationProperties *ApplicationProperties) fillMissingWithDefaultValues() {
	defaultProperties := getDefaultProperties()
	if applicationProperties.ApplicationServer == "" {
		applicationProperties.ApplicationServer = defaultProperties.ApplicationServer
	}
	if applicationProperties.ExasolHost == "" {
		applicationProperties.ExasolHost = defaultProperties.ExasolHost
	}
	if applicationProperties.ExasolPort == 0 {
		applicationProperties.ExasolPort = defaultProperties.ExasolPort
	}
	if applicationProperties.ExasolWebsocketAPIVersion == 0 {
		applicationProperties.ExasolWebsocketAPIVersion = defaultProperties.ExasolWebsocketAPIVersion
	}
}

func (applicationProperties *ApplicationProperties) validate() error {
	if applicationProperties.ExasolUser == "" && applicationProperties.ExasolPassword == "" {
		return fmt.Errorf(error_reporting_go.ExaError("E-ERA-8").
			Message("exasol username and password are missing in properties.").
			Mitigation("please specify an Exasol username and password via properties.").String())
	} else if applicationProperties.ExasolUser == "" {
		return fmt.Errorf(error_reporting_go.ExaError("E-ERA-9").
			Message("exasol username is missing in properties.").
			Mitigation("please specify an Exasol username via properties.").String())
	} else if applicationProperties.ExasolPassword == "" {
		return fmt.Errorf(error_reporting_go.ExaError("E-ERA-10").
			Message("exasol password is missing in properties.").
			Mitigation("please specify an Exasol password via properties.").String())
	} else {
		return nil
	}
}

func getDefaultProperties() *ApplicationProperties {
	return &ApplicationProperties{
		ApplicationServer:         "0.0.0.0:8080",
		ExasolHost:                "localhost",
		ExasolPort:                8563,
		ExasolWebsocketAPIVersion: 2,
	}
}
