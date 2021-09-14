package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	exasolrestapi "main/cmd/exasol-rest-api"
)

func main() {
	router := gin.Default()
	router.GET("/api/v1/query/:query", exasolrestapi.Query)
	appProperties := readApplicationProperties()
	err := router.Run(appProperties.Server)
	if err != nil {
		exasolrestapi.ErrorLogger.Printf("error starting API server: %s", err)
	}
}

func readApplicationProperties() applicationProperties {
	var appProperties applicationProperties
	var propertiesAsInterface interface{} = appProperties
	exasolrestapi.GetPropertiesFromFile("application-properties.yml", &propertiesAsInterface)
	err := mapstructure.Decode(propertiesAsInterface, &appProperties)
	if err != nil {
		exasolrestapi.ErrorLogger.Printf("error reading application properties: %s", err)
	}
	return appProperties
}

type applicationProperties struct {
	Server string
}
