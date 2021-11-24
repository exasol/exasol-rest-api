package exasol_rest_api

import (
	"github.com/exasol/error-reporting-go"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	_ "main/doc/swagger" // importing Swagger-generated documentation
	"net/http"
	"time"
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

func (application *Application) Auth(context *gin.Context) {
	err := application.Authorizer.Authorize(context.Request)
	if err != nil {
		context.JSON(http.StatusForbidden, APIBaseResponse{Status: "error", Exception: err.Error()})
		context.Abort()
	}
}

// AddEndpoints adds endpoints to the REST API.
func AddEndpoints(router *gin.Engine, application Application) {
	rate := limiter.Rate{Period: 1 * time.Minute, Limit: 30}
	rateLimiterMiddleware := mgin.NewMiddleware(limiter.New(memory.NewStore(), rate))

	router.ForwardedByClientIP = true
	router.Use(rateLimiterMiddleware)
	if application.Properties.APIAuth == 1 {
		addEndpointsWithAuth(router, application)
	} else {
		addEndpointsWithoutAuth(router, application)
	}
}

func addEndpointsWithAuth(router *gin.Engine, application Application) {
	router.GET("/api/v1/query/:query", application.Auth, application.Query)
	router.GET("/api/v1/tables", application.Auth, application.GetTables)
	router.GET("/api/v1/rows", application.Auth, application.GetRows)
	router.POST("/api/v1/row", application.Auth, application.InsertRow)
	router.DELETE("/api/v1/rows", application.Auth, application.DeleteRows)
	router.PUT("/api/v1/rows", application.Auth, application.UpdateRows)
	router.POST("/api/v1/statement", application.Auth, application.ExecuteStatement)
}

func addEndpointsWithoutAuth(router *gin.Engine, application Application) {
	router.GET("/api/v1/query/:query", application.Query)
	router.GET("/api/v1/tables", application.GetTables)
	router.GET("/api/v1/rows", application.GetRows)
	router.POST("/api/v1/row", application.InsertRow)
	router.DELETE("/api/v1/rows", application.DeleteRows)
	router.PUT("/api/v1/rows", application.UpdateRows)
	router.POST("/api/v1/statement", application.ExecuteStatement)
}
