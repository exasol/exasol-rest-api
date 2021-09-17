package exasol_rest_api

import (
	"github.com/gin-gonic/gin"
)

func Run() {
	applicationProperties := getApplicationProperties("APPLICATION_PROPERTIES_PATH")
	application := Application{
		Properties: applicationProperties,
	}
	router := gin.Default()
	router.GET("/api/v1/query/:query", application.Query)
	err := router.Run(applicationProperties.ApplicationServer)
	if err != nil {
		errorLogger.Printf("error starting API server: %s", err)
	}
}
