package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/api/v1/query/:query", query)
	appProperties := readApplicationProperties()
	err := router.Run(appProperties.Server)
	if err != nil {
		errorLogger.Printf("error starting API server: %s", err)
	}
}

func readApplicationProperties() applicationProperties {
	var appProperties applicationProperties
	var propertiesAsInterface interface{} = appProperties
	GetPropertiesFromFile("application-properties.yml", &propertiesAsInterface)
	err := mapstructure.Decode(propertiesAsInterface, &appProperties)
	if err != nil {
		errorLogger.Printf("error reading an application properties: %s", err)
	}
	return appProperties
}

func query(context *gin.Context) {
	response, err := Query(context.Param("query"))
	if err != nil {
		errorLogger.Printf("error during querying Exasol: %s", err)
	} else {
		context.IndentedJSON(http.StatusOK, response)
	}
}

type applicationProperties struct {
	Server string
}
