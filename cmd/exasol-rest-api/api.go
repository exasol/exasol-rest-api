package exasol_rest_api

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"os"
)

type applicationProperties struct {
	Server string
}

func Run() {
	appProperties := readApplicationProperties()
	connectionPropertiesPathKey := "CONNECTION_PROPERTIES_PATH"
	propertiesPath := os.Getenv(connectionPropertiesPathKey)
	if propertiesPath == "" {
		panic("runtime error: missing environment variable: " + connectionPropertiesPathKey)
	}
	properties := readConnectionProperties(connectionPropertiesPathKey)

	application := Application{
		ConnProperties: properties,
	}

	router := gin.Default()
	router.GET("/api/v1/query/:query", application.Query)

	err := router.Run(appProperties.Server)
	if err != nil {
		errorLogger.Printf("error starting API server: %s", err)
	}
}

func readApplicationProperties() applicationProperties {
	var appProperties applicationProperties
	var propertiesAsInterface interface{} = appProperties
	err := getPropertiesFromFile("application-properties.yml", &propertiesAsInterface)
	if err != nil {
		errorLogger.Printf("cannot extract application properties from a file: %s", err)
	}
	err = mapstructure.Decode(propertiesAsInterface, &appProperties)
	if err != nil {
		errorLogger.Printf("error reading application properties: %s", err)
	}
	return appProperties
}

func readConnectionProperties(propertiesPath string) *ConnectionProperties {
	var properties ConnectionProperties
	var propertiesAsInterface interface{} = properties
	err := getPropertiesFromFile(propertiesPath, &propertiesAsInterface)
	if err != nil {
		return nil
	}
	err = mapstructure.Decode(propertiesAsInterface, &properties)
	if err != nil {
		errorLogger.Printf("error reading websocketConnection properties: %s", err)
	}
	return createConnectionProperties(properties)
}
