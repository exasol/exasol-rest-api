/*
Package exasol_rest_api contains Exasol REST API logic.
*/
package exasol_rest_api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	exaerror "github.com/exasol/error-reporting-go"
	"github.com/exasol/exasol-driver-go"
	"github.com/gin-gonic/gin"
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
// [impl->dsn~execute-query-endpoint~1]
// [impl->dsn~execute-query-request-parameters~1]
func (application *Application) Query(context *gin.Context) {
	context.JSON(application.handleRequest(ConvertToGetRowsResponse, context.Param("query")))
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
// [impl->dsn~execute-statement-endpoint~1]
func (application *Application) ExecuteStatement(context *gin.Context) {
	var request ExecuteStatementRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(err))
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(validationError))
	} else {
		context.JSON(application.handleStatementRequest(request.GetStatement()))
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
// [impl->dsn~get-tables-endpoint~1]
func (application *Application) GetTables(context *gin.Context) {
	statement := "SELECT TABLE_SCHEMA, TABLE_NAME FROM EXA_ALL_TABLES"
	context.JSON(application.handleRequest(ConvertToGetTablesResponse, statement))
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
// [impl->dsn~insert-row-endpoint~1]
func (application *Application) InsertRow(context *gin.Context) {
	var request InsertRowRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(err))
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(validationError))
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		columnNames, values, err := request.GetRow()
		if err != nil {
			context.JSON(http.StatusBadRequest, apiErrorResponse(err))
		} else {
			statement := "INSERT INTO " + schemaName + "." + tableName + " (" + columnNames + ") VALUES (" + values + ")"
			context.JSON(application.handleStatementRequest(statement))
		}
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
// [impl->dsn~delete-rows-endpoint~1]
func (application *Application) DeleteRows(context *gin.Context) {
	var request RowsRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(err))
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(validationError))
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		condition, err := request.GetCondition()
		if err != nil {
			context.JSON(http.StatusBadRequest, apiErrorResponse(err))
		} else {
			statement := "DELETE FROM " + schemaName + "." + tableName + " WHERE " + condition
			context.JSON(application.handleStatementRequest(statement))
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
// [impl->dsn~update-rows-endpoint~1]
func (application *Application) UpdateRows(context *gin.Context) {
	var request UpdateRowsRequest
	err := context.BindJSON(&request)
	validationError := request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(err))
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(validationError))
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		valuesToUpdate, valuesError := request.GetValuesToUpdate()
		condition, conditionError := request.GetCondition()
		if valuesError != nil {
			context.JSON(http.StatusBadRequest, apiErrorResponse(valuesError))
		} else if conditionError != nil {
			context.JSON(http.StatusBadRequest, apiErrorResponse(conditionError))
		} else {
			statement := "UPDATE " + schemaName + "." + tableName + " SET " + valuesToUpdate + " WHERE " + condition
			context.JSON(application.handleStatementRequest(statement))
		}
	}
}

// @Summary GetRows from a table based on a condition
// @Description get zero or more rows from a table providing a WHERE condition
// @Produce  json
// @Security ApiKeyAuth
// @Param schemaName query string true "Exasol schema name"
// @Param tableName query string true "Exasol table name"
// @Param columnName query string false "Exasol column name for WHERE clause"
// @Param comparisonPredicate query string false "Comparison predicate for WHERE clause"
// @Param value query string false "Value of the specified Exasol column"
// @Param valueType query string false "Type of the value: string, bool, int or float"
// @Success 200 {string} status and response
// @Failure 400 {object} APIBaseResponse
// @Failure 403 {object} APIBaseResponse
// @Router /rows [get]
// [impl->dsn~get-rows-endpoint~1]
// [impl->dsn~get-rows-request-parameters~1]
func (application *Application) GetRows(context *gin.Context) {
	request, err := buildGetRowsRequest(context)
	validationError := request.ValidateWithOptionalCondition()
	if err != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(err))
	} else if validationError != nil {
		context.JSON(http.StatusBadRequest, apiErrorResponse(validationError))
	} else {
		schemaName := request.GetSchemaName()
		tableName := request.GetTableName()
		if !request.HasWhereClause() {
			statement := "SELECT * FROM " + schemaName + "." + tableName
			context.JSON(application.handleRequest(ConvertToGetRowsResponse, statement))
		} else {
			condition, conditionError := request.GetCondition()
			if conditionError != nil {
				context.JSON(http.StatusBadRequest, apiErrorResponse(conditionError))
			} else {
				statement := "SELECT * FROM " + schemaName + "." + tableName + " WHERE " + condition
				context.JSON(application.handleRequest(ConvertToGetRowsResponse, statement))
			}
		}
	}
}

func buildGetRowsRequest(context *gin.Context) (RowsRequest, error) {
	valueType := context.Query("valueType")
	value := context.Query("value")
	columnName := context.Query("columnName")
	comparisonPredicate := context.Query("comparisonPredicate")

	if valueType == "" && value == "" && columnName == "" && comparisonPredicate == "" {
		return RowsRequest{
			SchemaName: context.Query("schemaName"),
			TableName:  context.Query("tableName"),
		}, nil
	} else if valueType != "" && value != "" && columnName != "" {
		return createRowsRequestWithCondition(context, valueType, value, columnName, comparisonPredicate)
	} else {
		return RowsRequest{}, exaerror.New("E-ERA-30").
			Message("incomplete condition in the request.").
			Mitigation("provide 'columnName', 'valueType' and 'value' for the condition or remove the condition")
	}
}

func createRowsRequestWithCondition(context *gin.Context, valueType string, value string, columnName string, comparisonPredicate string) (RowsRequest, error) {
	renderedValue, err := getRenderedValue(context, valueType, value)
	if err != nil {
		return RowsRequest{}, err
	} else {
		return RowsRequest{
			SchemaName: context.Query("schemaName"),
			TableName:  context.Query("tableName"),
			WhereCondition: Condition{
				CellValue: Value{
					Value:      renderedValue,
					ColumnName: columnName,
				},
				ComparisonPredicate: comparisonPredicate,
			},
		}, nil
	}
}

func getRenderedValue(context *gin.Context, valueType string, value string) (interface{}, error) {
	if valueType != "" && value != "" {
		whereConditionValue, err := getValueByType(valueType, value)
		if err != nil {
			return nil, exaerror.New("E-ERA-28").
				Message("cannot decode value {{value}} with the provided value type {{value type}}: {{error}}").
				Parameter("value", context.Query("value")).
				Parameter("value type", context.Query("valueType")).
				Parameter("error", err.Error())
		} else {
			return whereConditionValue, nil
		}
	} else {
		return "", nil
	}
}

// [impl->dsn~execute-query-headers~1]
// [impl->dsn~get-tables-headers~1]
// [impl->dsn~insert-row-headers~1]
// [impl->dsn~delete-rows-headers~1]
// [impl->dsn~get-rows-headers~1]
// [impl->dsn~update-rows-headers~1]
// [impl->dsn~execute-statement-headers~1]
func (application *Application) handleRequest(convert func(toConvert *sql.Rows) (interface{}, error), query string) (int, interface{}) {
	connection, err := application.openConnection()

	if err != nil {
		wrappedError := exaerror.New("E-ERA-2").
			Message("error while opening a connection with Exasol: {{error|uq}}").
			Parameter("error", err.Error())
		return http.StatusInternalServerError, apiErrorResponse(wrappedError)
	}
	defer func() {
		err := connection.Close()
		if err != nil {
			errorLogger.Print("error closing connection: %w", err)
		}
	}()

	rows, err := connection.Query(query)
	if err != nil {
		wrappedError := exaerror.New("E-ERA-3").Message("error while executing query {{query}}: {{error|uq}}").
			Parameter("query", query).
			Parameter("error", err.Error())
		// Return 200 OK when query fails for backwards compatibility
		return http.StatusOK, apiErrorResponse(wrappedError)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			errorLogger.Print("error closing result set: %w", err)
		}
	}()

	convertedResponse, err := convert(rows)
	if err != nil {
		return http.StatusBadRequest, apiErrorResponse(err)
	} else {
		return http.StatusOK, convertedResponse
	}
}

func (application *Application) handleStatementRequest(statement string) (int, interface{}) {
	connection, err := application.openConnection()
	if err != nil {
		wrappedError := exaerror.New("E-ERA-2").
			Message("error while opening a connection with Exasol: {{error|uq}}").
			Parameter("error", err.Error())
		return http.StatusInternalServerError, apiErrorResponse(wrappedError)
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			errorLogger.Print("error closing connection: %w", err)
		}
	}()
	_, err = connection.Exec(statement)
	if err != nil {
		wrappedError := exaerror.New("E-ERA-31").Message("error while executing statement {{statement}}: {{error|uq}}").
			Parameter("statement", statement).
			Parameter("error", err.Error())
		// Return 200 OK when statement fails for backwards compatibility
		return http.StatusOK, apiErrorResponse(wrappedError)
	}
	return http.StatusOK, apiOkResponse()
}

func getValueByType(valueType string, valueAsString string) (interface{}, error) {
	switch valueType {
	case "string":
		return valueAsString, nil
	case "bool":
		return strconv.ParseBool(valueAsString)
	case "int":
		return strconv.Atoi(valueAsString)
	case "float":
		return strconv.ParseFloat(valueType, 64)
	default:
		return "", errors.New("unsupported value type: " + valueType)
	}
}

// [impl->dsn~communicate-with-database~2]
func (application *Application) openConnection() (*sql.DB, error) {
	props := application.Properties
	database, err := sql.Open("exasol", exasol.NewConfig(props.ExasolUser, props.ExasolPassword).
		Host(props.ExasolHost).
		Port(props.ExasolPort).
		Compression(true).
		Encryption(true). // Deactivating encryption not supported any more
		ValidateServerCertificate(props.ExasolValidateServerCertificate != "false").
		CertificateFingerprint(props.ExasolCertificateFingerprint).
		Autocommit(true).
		ClientName("Exasol REST API").
		String())

	if err != nil {
		return nil, err
	}
	// Verify that connection works to avoid later errors when executing queries
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	return database, nil
}
