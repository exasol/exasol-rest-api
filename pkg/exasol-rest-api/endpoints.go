/*
Package exasol_rest_api contains Exasol REST API logic.
*/
package exasol_rest_api

import (
	"errors"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
// @Success 200 {object} GetTablesResponse
// @Failure 400 {object} GetTablesResponse
// @Failure 403 {object} GetTablesResponse
// @Router /tables [get]
func (application *Application) GetTables(context *gin.Context) {
	statement := "SELECT TABLE_SCHEMA, TABLE_NAME FROM EXA_ALL_TABLES"
	err := application.Authorizer.Authorize(context.Request)
	if err != nil {
		context.JSON(http.StatusForbidden, GetTablesResponse{Status: "error", Exception: err.Error()})
	} else {
		application.handleGetTablesRequest(context, statement)
	}
}

func (application *Application) handleGetTablesRequest(context *gin.Context, statement string) {
	response, err := application.queryExasol(statement)
	if err != nil {
		context.JSON(http.StatusBadRequest, GetTablesResponse{Status: "error", Exception: err.Error()})
	} else {
		convertedResponse, err := ConvertToGetTablesResponse(response)
		if err != nil {
			context.JSON(http.StatusBadRequest, GetTablesResponse{Status: "error", Exception: err.Error()})
		} else {
			context.JSON(http.StatusOK, convertedResponse)
		}
	}
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
// @Param request-body body RowsRequest true "Request body"
// @Success 200 {string} status and response
// @Failure 400 {string} error code and error message
// @Failure 403 {string} error code and error message
// @Router /rows [delete]
func (application *Application) DeleteRows(context *gin.Context) {
	var request RowsRequest
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

// @Summary GetRows from a table based on a condition
// @Description get zero or more rows from a table providing a WHERE condition
// @Produce  json
// @Security ApiKeyAuth
// @Param schemaName query string true "Exasol schema name"
// @Param tableName query string true "Exasol table name"
// @Param columnName query string true "Exasol column name for WHERE clause"
// @Param comparisonPredicate query string true "Comparison predicate for WHERE clause"
// @Param value query string true "Value of the specified Exasol column"
// @Param valueType query string true "Type of the value: string, bool, int or float"
// @Success 200 {string} status and response
// @Failure 400 {string} error code and error message
// @Failure 403 {string} error code and error message
// @Router /rows [get]
func (application *Application) GetRows(context *gin.Context) {
	value, err := getValueByType(context.Query("valueType"), context.Query("value"))
	request := RowsRequest{
		SchemaName: context.Query("schemaName"),
		TableName:  context.Query("tableName"),
		WhereCondition: Condition{
			CellValue: Value{
				Value:      value,
				ColumnName: context.Query("columnName"),
			},
			ComparisonPredicate: context.Query("comparisonPredicate"),
		},
	}
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest,
			GetTablesResponse{Status: "error", Exception: error_reporting_go.ExaError("E-ERA-28").
				Message("cannot decode value {{value}} with the provided value type {{value type}}: {{error}}").
				Parameter("value", context.Query("value")).
				Parameter("value type", context.Query("valueType")).
				Parameter("error", err.Error()).String()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, GetTablesResponse{Status: "error", Exception: validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		condition, conditionError := request.GetCondition()
		if conditionError != nil {
			context.JSON(http.StatusBadRequest, GetTablesResponse{Status: "error", Exception: conditionError.Error()})
		} else {
			statement := "SELECT * FROM " + schemaName + "." + tableName + " WHERE " + condition
			application.executeGetRowsStatement(context, statement)
		}
	}
}

func (application *Application) executeGetRowsStatement(context *gin.Context, statement string) {
	err := application.Authorizer.Authorize(context.Request)
	if err != nil {
		context.JSON(http.StatusForbidden, GetTablesResponse{Status: "error", Exception: err.Error()})
	} else {
		application.handleGetRowsRequest(context, statement)
	}
}

func (application *Application) handleGetRowsRequest(context *gin.Context, statement string) {
	response, err := application.queryExasol(statement)
	if err != nil {
		context.JSON(http.StatusBadRequest, GetTablesResponse{Status: "error", Exception: err.Error()})
	} else {
		convertedResponse, err := ConvertToGetRowsResponse(response)
		if err != nil {
			context.JSON(http.StatusBadRequest, GetTablesResponse{Status: "error", Exception: err.Error()})
		} else {
			context.JSON(http.StatusOK, convertedResponse)
		}
	}
}

func getValueByType(valueType string, valueAsString string) (interface{}, error) {
	if valueType == "string" {
		return valueAsString, nil
	} else if valueType == "bool" {
		return strconv.ParseBool(valueAsString)
	} else if valueType == "int" {
		return strconv.Atoi(valueAsString)
	} else if valueType == "float" {
		return strconv.ParseFloat(valueType, 64)
	} else {
		return "", errors.New("unsupported value type: " + valueType)
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
		return nil, error_reporting_go.ExaError("E-ERA-3").Message("error while executing a query {{query}}: {{error|uq}}").
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
