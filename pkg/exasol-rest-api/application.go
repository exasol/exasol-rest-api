package exasol_rest_api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "main/doc/swagger"
)

// Run starts the REST API service.
func Run() {
	applicationProperties := GetApplicationProperties("APPLICATION_PROPERTIES_PATH")
	application := Application{
		Properties: applicationProperties,
	}
	router := gin.Default()
	router.GET("/api/v1/query/:query", application.Query)
	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	err := router.Run(applicationProperties.ApplicationServer)
	if err != nil {
		panic("error starting API server: " + err.Error())
	}
}
