package exasol_rest_api

import (
	"errors"
	"os"
)

//ApplicationProperties for Exasol REST API service.
type ApplicationProperties struct {
	ApplicationServer         string `yaml:"server-address"`
	ExasolUser                string `yaml:"exasol-user"`
	ExasolPassword            string `yaml:"exasol-password"`
	ExasolHost                string `yaml:"exasol-host"`
	ExasolPort                int    `yaml:"exasol-port"`
	ExasolWebsocketApiVersion int    `yaml:"exasol-websocket-api-version"`
	Encryption                bool   `yaml:"encryption"`
	UseTLS                    bool   `yaml:"use-tls"`
}

//GetApplicationProperties creates an application properties.
func GetApplicationProperties(applicationPropertiesPathKey string) *ApplicationProperties {
	propertiesPath := os.Getenv(applicationPropertiesPathKey)
	if propertiesPath == "" {
		panic("runtime error: missing environment variable: " + applicationPropertiesPathKey)
	}
	properties, err := readApplicationProperties(propertiesPath)
	if err != nil {
		panic("runtime error: application properties are missing or incorrect. " + err.Error())
	}
	return properties
}

func readApplicationProperties(propertiesFilePath string) (*ApplicationProperties, error) {
	var properties ApplicationProperties
	err := getPropertiesFromFile(propertiesFilePath, &properties)
	if err != nil {
		return nil, err
	}
	err = properties.initializeProperties()
	if err != nil {
		return nil, err
	} else {
		return &properties, nil
	}
}
func (applicationProperties *ApplicationProperties) initializeProperties() error {
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
	if applicationProperties.ExasolWebsocketApiVersion == 0 {
		applicationProperties.ExasolWebsocketApiVersion = defaultProperties.ExasolWebsocketApiVersion
	}
	return applicationProperties.validateExasolUser()
}

func (applicationProperties *ApplicationProperties) validateExasolUser() error {
	if applicationProperties.ExasolUser == "" && applicationProperties.ExasolPassword == "" {
		return errors.New("exasol username and password are missing in properties")
	} else if applicationProperties.ExasolUser == "" {
		return errors.New("exasol username is missing in properties")
	} else if applicationProperties.ExasolPassword == "" {
		return errors.New("exasol password is missing in properties")
	} else {
		return nil
	}
}

func getDefaultProperties() *ApplicationProperties {
	return &ApplicationProperties{
		ApplicationServer:         "localhost:8080",
		ExasolHost:                "localhost",
		ExasolPort:                8563,
		ExasolWebsocketApiVersion: 2,
	}
}
