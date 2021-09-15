package exasol_rest_api

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type applicationProperties struct {
	Server string
}

func Run() {
	appProperties := readApplicationProperties()
	router := gin.Default()

	router.GET("/api/v1/query/:query", Query)

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
