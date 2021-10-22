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

// @Summary Query the Exasol database.
// @Description provide a query and get a result set
// @Produce  json
// @Security ApiKeyAuth
// @Param   query     path    string     true        "SELECT query"
// @Success 200 {string} status and result set
// @Failure 400 {string} error code and error message
// @Failure 403 {string} error code and error message
// @Router /query/{query} [get]
func (application *Application) Query(context *gin.Context) {
	application.executeStatement(context, context.Param("query"))
}

// @Summary GetTables that are available for the user.
// @Description get a list of all available tables
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} status and result set
// @Failure 400 {string} error code and error message
// @Failure 403 {string} error code and error message
// @Router /tables [get]
func (application *Application) GetTables(context *gin.Context) {
	statement := "SELECT * FROM EXA_USER_TABLES"
	application.executeStatement(context, statement)
}

// @Summary InsertRow to a table.
// @Description insert a single row into an Exasol table
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param request-body body InsertRowRequest true "Request body"
// @Success 200 {string} status and response
// @Failure 400 {string} error code and error message
// @Failure 403 {string} error code and error message
// @Router /row [post]
func (application *Application) InsertRow(context *gin.Context) {
	var request InsertRowRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		columnNames, values, _ := request.GetRow()
		statement := "INSERT INTO " + schemaName + "." + tableName + " (" + columnNames + ") VALUES (" + values + ")"
		application.executeStatement(context, statement)
	}
}

// @Summary DeleteRows from a table based on a condition
// @Description delete zero or more rows from a table providing a WHERE condition
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param request-body body DeleteRowsRequest true "Request body"
// @Success 200 {string} status and response
// @Failure 400 {string} error code and error message
// @Failure 403 {string} error code and error message
// @Router /rows [delete]
func (application *Application) DeleteRows(context *gin.Context) {
	var request DeleteRowsRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		condition, err := request.GetCondition()
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		} else {
			statement := "DELETE FROM " + schemaName + "." + tableName + " WHERE " + condition
			application.executeStatement(context, statement)
		}
	}
}

// @Summary UpdateRows in a table based on a condition
// @Description update zero or more row in a table based on a condition
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param request-body body UpdateRowsRequest true "Request body"
// @Success 200 {string} status and response
// @Failure 400 {string} error code and error message
// @Failure 403 {string} error code and error message
// @Router /rows [put]
func (application *Application) UpdateRows(context *gin.Context) {
	var request UpdateRowsRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		valuesToUpdate, valuesError := request.GetValuesToUpdate()
		condition, conditionError := request.GetCondition()
		if valuesError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Error": valuesError.Error()})
		} else if conditionError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Error": conditionError.Error()})
		} else {
			statement := "UPDATE " + schemaName + "." + tableName + " SET " + valuesToUpdate + " WHERE " + condition
			application.executeStatement(context, statement)
		}
	}
}

func (application *Application) executeStatement(context *gin.Context, query string) {
	err := application.Authorizer.Authorize(context.Request)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"Error": err.Error()})
	} else {
		response, err := application.queryExasol(query)
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
