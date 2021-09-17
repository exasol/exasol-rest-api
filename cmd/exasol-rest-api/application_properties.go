package exasol_rest_api

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"os"
)

type ApplicationProperties struct {
	ApplicationServer         string
	ExasolUser                string
	ExasolPassword            string
	ExasolHost                string
	ExasolPort                int
	ExasolWebsocketApiVersion int
	Encryption                bool
	UseTLS                    bool
}

func getApplicationProperties(applicationPropertiesPathKey string) *ApplicationProperties {
	propertiesPath := os.Getenv(applicationPropertiesPathKey)
	if propertiesPath == "" {
		panic("runtime error: missing environment variable: " + applicationPropertiesPathKey)
	}
	properties, err := readApplicationProperties(propertiesPath)
	if err != nil {
		errorLogger.Printf("runtime error: cannot read application properties. %s", err)
		panic("cannot start application without properties")
	}
	return properties
}

func readApplicationProperties(propertiesFilePath string) (*ApplicationProperties, error) {
	var properties ApplicationProperties
	var propertiesAsInterface interface{} = properties
	err := getPropertiesFromFile(propertiesFilePath, &propertiesAsInterface)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(propertiesAsInterface, &properties)
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
