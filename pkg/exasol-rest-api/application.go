package exasol_rest_api

import (
	"github.com/exasol/error-reporting-go"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "main/doc/swagger" // importing Swagger-generated documentation
)

// Run starts the REST API service.
func Run() {
	applicationProperties := GetApplicationProperties("APPLICATION_PROPERTIES_PATH")
	application := Application{
		Properties: applicationProperties,
	}
	router := gin.Default()
	swaggerURL := ginSwagger.URL("/swagger/doc.json")

	router.GET("/api/v1/query/:query", application.Query)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	err := router.Run(applicationProperties.ApplicationServer)

	if err != nil {
		panic(error_reporting_go.ExaError("E-ERA-1").Message("error starting API server: {{error}}").
			Parameter("error", err.Error()).String())
	}
}
