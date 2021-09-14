package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	exasol_rest_api "main/cmd/exasol-rest-api"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/api/v1/query/:query", query)
	appProperties := readApplicationProperties()
	err := router.Run(appProperties.Server)
	if err != nil {
		exasol_rest_api.ErrorLogger.Printf("error starting API server: %s", err)
	}
}

func readApplicationProperties() applicationProperties {
	var appProperties applicationProperties
	var propertiesAsInterface interface{} = appProperties
	exasol_rest_api.GetPropertiesFromFile("application-properties.yml", &propertiesAsInterface)
	err := mapstructure.Decode(propertiesAsInterface, &appProperties)
	if err != nil {
		exasol_rest_api.ErrorLogger.Printf("error reading application properties: %s", err)
	}
	return appProperties
}

func query(context *gin.Context) {
	response, err := exasol_rest_api.Query(context.Param("query"))
	if err != nil {
		exasol_rest_api.ErrorLogger.Printf("error during querying Exasol: %s", err)
	} else {
		context.IndentedJSON(http.StatusOK, response)
	}
}

type applicationProperties struct {
	Server string
}
