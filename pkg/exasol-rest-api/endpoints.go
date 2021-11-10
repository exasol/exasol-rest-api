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
// @Success 200 {string} status and response
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /query/{query} [get]
func (application *Application) Query(context *gin.Context) {
	context.JSON(application.handleRequest(ConvertToGetRowsResponse, context.Request,
		context.Param("query")))
}

// @Summary ExecuteStatement on the Exasol database.
// @Description execute a statement without a result set
// @Produce  json
// @Security ApiKeyAuth
// @Param request-body body ExecuteStatementRequest true "Request body"
// @Success 200 {string} APIBaseResponse
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /statement [post]
func (application *Application) ExecuteStatement(context *gin.Context) {
	var request ExecuteStatementRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: validationError.Error()})
	} else {
		context.JSON(application.handleRequest(ConvertToBaseResponse, context.Request, request.GetStatement()))
	}
}

// @Summary GetTables that are available for the user.
// @Description get a list of all available tables
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} GetTablesResponse
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /tables [get]
func (application *Application) GetTables(context *gin.Context) {
	statement := "SELECT TABLE_SCHEMA, TABLE_NAME FROM EXA_ALL_TABLES"
	context.JSON(application.handleRequest(ConvertToGetTablesResponse, context.Request, statement))
}

// @Summary InsertRow to a table.
// @Description insert a single row into an Exasol table
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param request-body body InsertRowRequest true "Request body"
// @Success 200 {object} APIBaseResponse
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /row [post]
func (application *Application) InsertRow(context *gin.Context) {
	var request InsertRowRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		columnNames, values, err := request.GetRow()
		if err != nil {
			context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()})
		}
		statement := "INSERT INTO " + schemaName + "." + tableName + " (" + columnNames + ") VALUES (" + values + ")"
		context.JSON(application.handleRequest(ConvertToBaseResponse, context.Request, statement))
	}
}

// @Summary DeleteRows from a table based on a condition
// @Description delete zero or more rows from a table providing a WHERE condition
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param request-body body RowsRequest true "Request body"
// @Success 200 {object} APIBaseResponse
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /rows [delete]
func (application *Application) DeleteRows(context *gin.Context) {
	var request RowsRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		condition, err := request.GetCondition()
		if err != nil {
			context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()})
		} else {
			statement := "DELETE FROM " + schemaName + "." + tableName + " WHERE " + condition
			context.JSON(application.handleRequest(ConvertToBaseResponse, context.Request, statement))
		}
	}
}

// @Summary UpdateRows in a table based on a condition
// @Description update zero or more row in a table based on a condition
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param request-body body UpdateRowsRequest true "Request body"
// @Success 200 {object} APIBaseResponse
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /rows [put]
func (application *Application) UpdateRows(context *gin.Context) {
	var request UpdateRowsRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		valuesToUpdate, valuesError := request.GetValuesToUpdate()
		condition, conditionError := request.GetCondition()
		if valuesError != nil {
			context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: valuesError.Error()})
		} else if conditionError != nil {
			context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: conditionError.Error()})
		} else {
			statement := "UPDATE " + schemaName + "." + tableName + " SET " + valuesToUpdate + " WHERE " + condition
			context.JSON(application.handleRequest(ConvertToBaseResponse, context.Request, statement))
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
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /rows [get]
func (application *Application) GetRows(context *gin.Context) {
	value, err := getValueByType(context.Query("valueType"), context.Query("value"))
	request := buildGetRowsRequest(context, value)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest,
			APIBaseResponse{Status: "error", Exception: error_reporting_go.ExaError("E-ERA-28").
				Message("cannot decode value {{value}} with the provided value type {{value type}}: {{error}}").
				Parameter("value", context.Query("value")).
				Parameter("value type", context.Query("valueType")).
				Parameter("error", err.Error()).String()})
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: validationError.Error()})
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		condition, conditionError := request.GetCondition()
		if conditionError != nil {
			context.JSON(http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: conditionError.Error()})
		} else {
			statement := "SELECT * FROM " + schemaName + "." + tableName + " WHERE " + condition
			context.JSON(application.handleRequest(ConvertToGetRowsResponse, context.Request, statement))
		}
	}
}

func buildGetRowsRequest(context *gin.Context, value interface{}) RowsRequest {
	return RowsRequest{
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
}

func (application *Application) handleRequest(convert func(toConvert []byte) (interface{}, error),
	request *http.Request, statement string) (int, interface{}) {
	err := application.Authorizer.Authorize(request)
	if err != nil {
		return http.StatusForbidden, APIBaseResponse{Status: "error", Exception: err.Error()}
	} else {
		return application.SendRequestToWebSocketsAPI(statement, convert)
	}
}

func (application *Application) SendRequestToWebSocketsAPI(statement string,
	convert func(toConvert []byte) (interface{}, error)) (int, interface{}) {
	response, err := application.queryExasol(statement)
	if err != nil {
		return http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()}
	} else {
		convertedResponse, err := convert(response)
		if err != nil {
			return http.StatusBadRequest, APIBaseResponse{Status: "error", Exception: err.Error()}
		} else {
			return http.StatusOK, convertedResponse
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
