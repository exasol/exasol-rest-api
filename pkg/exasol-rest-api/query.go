/*
Package exasol_rest_api contains Exasol REST API logic.
*/
package exasol_rest_api

import (
	error_reporting_go "github.com/exasol/error-reporting-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Application represents the REST API service.
type Application struct {
	Properties *ApplicationProperties
	Authorizer Authorizer
}

// @Summary Query the Exasol databse.
// @Description provide a query and get a result set
// @Accept  json
// @Produce  json
// @Param   query     path    string     true        "SELECT query"
// @Success 200 {string} status and result set
// @Failure 400 {string} error code and error message
// @Router /query/{query} [get]
func (application *Application) Query(context *gin.Context) {
	err := application.Authorizer.Authorize(context.Request)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"Error": err.Error()})
	} else {
		response, err := application.queryExasol(context.Param("query"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		} else {
			context.Data(http.StatusOK, "application/json", response)
		}
	}
}

func (application *Application) queryExasol(query string) ([]byte, error) {
	connection, err := application.openConnection()
	if err != nil {
		return nil, error_reporting_go.ExaError("E-ERA-2").
			Message("error while opening a connection with Exasol: {{error|uq}}").
			Parameter("error", err.Error())
	}

	defer connection.close()

	response, err := connection.executeQuery(query)
	if err != nil {
		return nil, error_reporting_go.ExaError("E-ERA-3").Message("error while executing a query: {{error|uq}}").
			Parameter("error", err.Error())
	}

	return response, nil
}

func (application *Application) openConnection() (*websocketConnection, error) {
	connection := &websocketConnection{
		connProperties: application.Properties,
	}

	err := connection.connect()
	if err != nil {
		return nil, err
	}

	err = connection.login()
	if err != nil {
		return nil, err
	}

	return connection, nil
}
