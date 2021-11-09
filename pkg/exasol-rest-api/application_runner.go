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
	applicationProperties := GetApplicationProperties()
	application := Application{
		Properties: applicationProperties,
		Authorizer: &TokenAuthorizer{
			AllowedTokens: CreateStringsSet(applicationProperties.APITokens),
		},
	}
	router := gin.Default()
	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	AddEndpoints(router, application)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	err := router.Run(applicationProperties.ApplicationServer)

	if err != nil {
		panic(error_reporting_go.ExaError("E-ERA-1").Message("error starting API server: {{error}}").
			Parameter("error", err.Error()).String())
	}
}

// AddEndpoints adds endpoints to the REST API.
func AddEndpoints(router *gin.Engine, application Application) {
	router.GET("/api/v1/query/:query", application.Query)
	router.GET("/api/v1/tables", application.GetTables)
	router.GET("/api/v1/rows", application.GetRows)
	router.POST("/api/v1/row", application.InsertRow)
	router.DELETE("/api/v1/rows", application.DeleteRows)
	router.PUT("/api/v1/rows", application.UpdateRows)
	router.POST("/api/v1/statement", application.ExecuteStatement)
}
