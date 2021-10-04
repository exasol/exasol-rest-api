package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"os"
)

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
func GetApplicationProperties(applicationPropertiesPathKey string) *ApplicationProperties {
	propertiesPath := os.Getenv(applicationPropertiesPathKey)
	if propertiesPath == "" {
		panic(error_reporting_go.ExaError("E-ERA-4").Message("missing environment variable: {{env}}.").
			Parameter("env", applicationPropertiesPathKey).
			Mitigation("please set the variable according to the user guide.").String())
	}

	properties, err := readApplicationProperties(propertiesPath)
	if err != nil {
		panic(error_reporting_go.ExaError("E-ERA-5").
			Message("application properties are missing or incorrect. {{error|uq}}").
			Parameter("error", err.Error()).String())
	}

	return properties
}

func readApplicationProperties(propertiesFilePath string) (*ApplicationProperties, error) {
	properties, err := getPropertiesFromFile(propertiesFilePath)
	if err != nil {
		return nil, error_reporting_go.ExaError("E-ERA-6").
			Message("cannot read properties from specified file: {{file path}}. {{error|uq}}").
			Parameter("file path", propertiesFilePath).Parameter("error", err.Error())
	}

	properties.fillMissingWithDefaultValues()

	err = properties.validate()
	if err != nil {
		return nil, error_reporting_go.ExaError("E-ERA-7").Message("properties file validation failed. {{error|uq}}").
			Parameter("error", err.Error())
	} else {
		return properties, nil
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
