package exasol_rest_api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/suite"

	testSetupAbstraction "github.com/exasol/exasol-test-setup-abstraction-server/go-client"
)

type IntegrationTestSuite struct {
	suite.Suite
	ctx                    context.Context
	exasolContainer        *testSetupAbstraction.TestSetupAbstraction
	defaultServiceUsername string
	defaultServicePassword string
	defaultAuthTokens      []string
	exasolPort             int
	exasolHost             string
	appProperties          *exasol_rest_api.ApplicationProperties
	connection             *sql.DB
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	if testing.Short() {
		suite.T().Skip()
	}
	suite.ctx = context.Background()
	suite.defaultServiceUsername = "api_service_account"
	suite.defaultServicePassword = "secret_password"
	suite.defaultAuthTokens = []string{"3J90XAv9loMIXzQdfYmtJrHAbopPsc", "OR6rq6KjWmhvGU770A9OTjpfH86nlk"}
	suite.exasolContainer = runExasolContainer()
	connectionInfo, err := suite.exasolContainer.GetConnectionInfo()
	onError(err)
	suite.exasolHost = connectionInfo.Host
	suite.exasolPort = connectionInfo.Port

	database, err := suite.exasolContainer.CreateConnectionWithConfig(true)
	onError(err)
	suite.connection = database
	createDefaultServiceUserWithAccess(suite.connection, suite.defaultServiceUsername, suite.defaultServicePassword)
}

func (suite *IntegrationTestSuite) startServer(application exasol_rest_api.Application) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	exasol_rest_api.AddEndpoints(router, application)
	suite.appProperties = application.Properties
	return router
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestQuery() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "SELECT * FROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"rows\":[{\"X\":15,\"Y\":\"test\"},{\"X\":10,\"Y\":\"test_2\"}],\"meta\":{\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100}}]}}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestQueryWithTypo() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "SELECTFROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-3: error while executing query 'SELECTFROM TEST_SCHEMA_1.TEST_TABLE': E-EGOD-11: execution failed with SQL error code '42000' and message 'syntax error, unexpected ",
	}
	suite.assertResponseBodyContains(&data, suite.sendQueryRequest(&data))
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestInsertNotAllowed() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "CREATE SCHEMA not_allowed_schema",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-3: error while executing query 'CREATE SCHEMA not_allowed_schema': E-EGOD-11: execution failed with SQL error code '42500' and message 'insufficient privileges for creating schema",
	}
	suite.assertResponseBodyContains(&data, suite.sendQueryRequest(&data))
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestExasolUserWithoutCreateSessionPrivilege() {
	username := "user_without_session_privilege"
	password := "secret"
	suite.createExasolUser(username, password)
	data := testData{
		server:         suite.createServerWithUser(username, password),
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusInternalServerError,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-2: error while opening a connection with Exasol: failed to login: E-EGOD-11: execution failed with SQL error code '08004' and message 'Connection exception - insufficient privileges: CREATE SESSION.'\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestCertificateValidationFailsWithoutFingerprint() {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                       suite.defaultAuthTokens,
		ExasolUser:                      suite.defaultServiceUsername,
		ExasolPassword:                  suite.defaultServicePassword,
		ExasolHost:                      suite.exasolHost,
		ExasolPort:                      suite.exasolPort,
		ExasolValidateServerCertificate: "true",
		ExasolCertificateFingerprint:    "",
	})

	data := testData{
		server:         server,
		query:          "select 1",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusInternalServerError,
		expectedBody:   `{"status":"error","exception":"E-ERA-2: error while opening a connection with Exasol: failed to connect to URL \"wss://localhost:32805\": tls: failed to verify certificate: x509: “exacluster.local” certificate is not standards compliant"}`,
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

// Get actual fingerprint of Exasol Docker container by extracting the fingerprint
// from the error message returned when connecting with an invalid fingerprint.
func (suite *IntegrationTestSuite) TestCertificateValidationSucceedsWithFingerprint() {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                       suite.defaultAuthTokens,
		ExasolUser:                      suite.defaultServiceUsername,
		ExasolPassword:                  suite.defaultServicePassword,
		ExasolHost:                      suite.exasolHost,
		ExasolPort:                      suite.exasolPort,
		ExasolValidateServerCertificate: "true",
		ExasolCertificateFingerprint:    suite.getExasolCertificateFingerprint(),
	})

	data := testData{
		server:         server,
		query:          "select 1",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   `{"status":"ok","rows":[[null,1]],"meta":{"columns":[{"name":"1","dataType":{"type":"DECIMAL","precision":1}}]}}`,
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) getExasolCertificateFingerprint() string {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                       suite.defaultAuthTokens,
		ExasolUser:                      suite.defaultServiceUsername,
		ExasolPassword:                  suite.defaultServicePassword,
		ExasolHost:                      suite.exasolHost,
		ExasolPort:                      suite.exasolPort,
		ExasolValidateServerCertificate: "true",
		ExasolCertificateFingerprint:    "invalidFingerprint",
	})

	data := testData{
		server:         server,
		query:          "ignored",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "",
	}
	response := suite.sendQueryRequest(&data)
	reg := regexp.MustCompile("the server's certificate fingerprint '([a-zA-Z0-9]+)' does not match the expected fingerprint 'invalidFingerprint'")
	submatches := reg.FindStringSubmatch(response.Body.String())
	actualFingerprint := submatches[1]
	if actualFingerprint == "" {
		suite.FailNowf("Expected response %q to match %q", response.Body.String(), reg)
	}
	return actualFingerprint
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestExasolUserWithWrongCredentials() {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                       suite.defaultAuthTokens,
		ExasolUser:                      "not_existing_user",
		ExasolPassword:                  "wrong_password",
		ExasolHost:                      suite.exasolHost,
		ExasolPort:                      suite.exasolPort,
		ExasolValidateServerCertificate: "false",
	})
	data := testData{
		server:         server,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusInternalServerError,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-2: error while opening a connection with Exasol: failed to login: E-EGOD-11: execution failed with SQL error code '08004' and message 'Connection exception - authentication failed.'\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestWrongExasolPort() {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:      suite.defaultAuthTokens,
		ExasolUser:     suite.defaultServiceUsername,
		ExasolPassword: suite.defaultServicePassword,
		ExasolHost:     suite.exasolHost,
		ExasolPort:     4321,
	})
	data := testData{
		server:         server,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusInternalServerError,
		expectedBody:   "failed to connect to URL \\\"wss://localhost:4321\\\":",
	}
	suite.assertResponseBodyContains(&data, suite.sendQueryRequest(&data))
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestUnauthorizedAccessToQuery() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "some query",
		authToken:      "OR6rq6KjWmhvGU770A9OTjpfH86nlkq",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-22: an authorization token is missing or wrong. please make sure you provided a valid token.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

// [itest->dsn~execute-query-endpoint~1]
// [itest->dsn~execute-query-response-body~2]
func (suite *IntegrationTestSuite) TestUnauthorizedAccessWithShortToken() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "some query",
		authToken:      "tooshort",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-23: an authorization token has invalid length: 8. please only use tokens with the length longer or equal to 30.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

// [itest->dsn~get-tables-endpoint~1]
// [itest->dsn~get-tables-response-body~1]
func (suite *IntegrationTestSuite) TestGetTables() {
	username := "GET_TABLES_USER"
	password := "secret"
	schemaName := "TEST_SCHEMA_GET_TABLES_1"
	columns := "C1 VARCHAR(100), C2 DECIMAL(5,0)"

	suite.creatSchemaAndTable(schemaName, "TEST_TABLE_1", columns)
	suite.creatSchemaAndTable(schemaName, "TEST_TABLE_2", columns)
	suite.createExasolUser(username, password)
	suite.grantToUser(username, "CREATE SESSION")
	suite.grantToUser(username, "SELECT ON SCHEMA "+schemaName)

	data := testData{
		server:         suite.createServerWithUser(username, password),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"tablesList\":[{\"tableName\":\"TEST_TABLE_1\",\"schemaName\":\"TEST_SCHEMA_GET_TABLES_1\"},{\"tableName\":\"TEST_TABLE_2\",\"schemaName\":\"TEST_SCHEMA_GET_TABLES_1\"}]}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetTables(&data))
}

// [itest->dsn~get-tables-endpoint~1]
// [itest->dsn~get-tables-response-body~1]
func (suite *IntegrationTestSuite) TestGetTablesWithZeroTables() {
	username := "USER_WITHOUT_OWNED_SCHEMA"
	password := "secret"
	suite.createExasolUser(username, password)
	suite.grantToUser(username, "CREATE SESSION")

	data := testData{
		server:         suite.createServerWithUser(username, password),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"tablesList\":[]}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetTables(&data))
}

// [itest->dsn~get-tables-endpoint~1]
// [itest->dsn~get-tables-response-body~1]
func (suite *IntegrationTestSuite) TestGetTablesUnauthorizedAccess() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "some query",
		authToken:      "OR6rq6KjWmhvGU770A9OTjpfH86nlkq",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-22: an authorization token is missing or wrong. please make sure you provided a valid token.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetTables(&data))
}

// [itest->dsn~insert-row-endpoint~1]
// [itest->dsn~insert-row-request-body~1]
// [itest->dsn~insert-row-response-body~1]
func (suite *IntegrationTestSuite) TestInsertRow() {
	username := "INSERT_ROW_USER"
	password := "secret"
	schemaName := "TEST_SCHEMA_INSERT_ROW_1"
	tableName := "ALL_DATA_TYPES"
	columns := "C1 VARCHAR(100), C2 VARCHAR(100) CHARACTER SET ASCII, C3 CHAR(10), C4 CHAR(10) CHARACTER SET ASCII, " +
		"C5 DECIMAL(5,0), C6 DECIMAL(6,3), C7 DOUBLE, C8 BOOLEAN, C9 DATE, C10 TIMESTAMP, " +
		"C11 TIMESTAMP WITH LOCAL TIME ZONE, C12 INTERVAL YEAR TO MONTH, C13 INTERVAL DAY TO SECOND, C14 GEOMETRY(3857)"

	suite.creatSchemaAndTable(schemaName, tableName, columns)
	suite.createExasolUser(username, password)
	suite.grantToUser(username, "CREATE SESSION")
	suite.grantToUser(username, "INSERT ON SCHEMA "+schemaName)

	data := testData{
		server:         suite.createServerWithUser(username, password),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\"}",
	}
	insertRowRequest := exasol_rest_api.InsertRowRequest{
		SchemaName: schemaName,
		TableName:  tableName,
		Row: []exasol_rest_api.Value{
			{ColumnName: "C1", Value: "Exa'sol"},
			{ColumnName: "C2", Value: "b"},
			{ColumnName: "C3", Value: "c"},
			{ColumnName: "C4", Value: "d"},
			{ColumnName: "C5", Value: 3},
			{ColumnName: "C6", Value: 123.456},
			{ColumnName: "C7", Value: 2.2},
			{ColumnName: "C8", Value: false},
			{ColumnName: "C9", Value: "2016-08-01"},
			{ColumnName: "C10", Value: "2016-08-01 23:12:01.000"},
			{ColumnName: "C11", Value: "2016-08-01 00:00:02.000"},
			{ColumnName: "C12", Value: "4-6"},
			{ColumnName: "C13", Value: "3 12:50:10.123"},
			{ColumnName: "C14", Value: "POINT(2 5)"},
		},
	}
	body, err := json.Marshal(insertRowRequest)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendInsertRow(&data, body))
	suite.assertInsertRowValuesInTable(schemaName, tableName)
}

// [itest->dsn~insert-row-endpoint~1]
// [itest->dsn~insert-row-request-body~1]
// [itest->dsn~insert-row-response-body~1]
func (suite *IntegrationTestSuite) assertInsertRowValuesInTable(schemaName string, tableName string) {
	rows, err := suite.connection.Query("SELECT * FROM " + schemaName + "." + tableName)
	onError(err)
	defer func() { onError(rows.Close()) }()
	rows.Next()

	var c1, c2, c3, c4, c9, c10, c11, c12, c13, c14 string
	var c5 int
	var c6, c7 float64
	var c8 bool

	err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8, &c9, &c10, &c11, &c12, &c13, &c14)
	onError(err)
	suite.Equal("Exa'sol", c1)
	suite.Equal("b", c2)
	suite.Equal("c         ", c3)
	suite.Equal("d         ", c4)
	suite.Equal(3, c5)
	suite.Equal(123.456, c6)
	suite.Equal(2.2, c7)
	suite.Equal(false, c8)
	suite.Equal("2016-08-01", c9)
	suite.Equal("2016-08-01 23:12:01.000000", c10)
	suite.Equal("2016-08-01 00:00:02.000000", c11)
	suite.Equal("+04-06", c12)
	suite.Equal("+03 12:50:10.123", c13)
	suite.Equal("POINT (2 5)", c14)
}

// [itest->dsn~insert-row-endpoint~1]
// [itest->dsn~insert-row-request-body~1]
// [itest->dsn~insert-row-response-body~1]
func (suite *IntegrationTestSuite) TestInsertRowAuthorizationError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "badToken",
		expectedStatus: http.StatusForbidden,
		expectedBody: "{\"status\":\"error\",\"exception\":\"E-ERA-23: an authorization token has invalid length: 8. " +
			"please only use tokens with the length longer or equal to 30.\"}",
	}
	insertRowRequest := exasol_rest_api.InsertRowRequest{
		SchemaName: "foo",
		TableName:  "bar",
		Row: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
	}
	body, err := json.Marshal(insertRowRequest)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendInsertRow(&data, body))
}

// [itest->dsn~insert-row-endpoint~1]
// [itest->dsn~insert-row-request-body~1]
// [itest->dsn~insert-row-response-body~1]
func (suite *IntegrationTestSuite) TestInsertRowMissingRequestParameter() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody: "{\"status\":\"error\",\"exception\":\"E-ERA-17: insert row request has some missing parameters. " +
			"Please specify schema name, table name and row\"}",
	}
	insertRowRequest := exasol_rest_api.InsertRowRequest{
		TableName: "bar",
		Row: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
	}
	body, err := json.Marshal(insertRowRequest)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendInsertRow(&data, body))
}

// [itest->dsn~delete-rows-endpoint~1]
// [itest->dsn~delete-rows-request-body~1]
// [itest->dsn~delete-rows-response-body~1]
func (suite *IntegrationTestSuite) TestDeleteRow() {
	username := "DELETE_ROWS_USER"
	password := "secret"
	schemaName := "TEST_SCHEMA_DELETE_ROWS_1"
	tableName := "TEST_TABLE"
	columns := "C1 VARCHAR(100), C2 DECIMAL(5,0)"

	suite.creatSchemaAndTable(schemaName, tableName, columns)
	suite.insertRowIntoTable(schemaName, tableName, "'row1', 1")
	suite.insertRowIntoTable(schemaName, tableName, "'row2', 2")
	suite.insertRowIntoTable(schemaName, tableName, "'row3', 2")
	suite.createExasolUser(username, password)
	suite.grantToUser(username, "CREATE SESSION")
	suite.grantToUser(username, "DELETE ON SCHEMA "+schemaName)

	data := testData{
		server:         suite.createServerWithUser(username, password),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\"}",
	}
	deleteRowsRequest := exasol_rest_api.RowsRequest{
		SchemaName: schemaName,
		TableName:  tableName,
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "C2",
				Value:      2,
			},
		},
	}
	body, err := json.Marshal(deleteRowsRequest)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendDeleteRows(&data, body))
	suite.assertTableHasOnlyOneRow(schemaName, tableName)
}

// [itest->dsn~delete-rows-endpoint~1]
// [itest->dsn~delete-rows-request-body~1]
// [itest->dsn~delete-rows-response-body~1]
func (suite *IntegrationTestSuite) TestDeleteRowsAuthorizationError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "12345678912345678912345678912345",
		expectedStatus: http.StatusForbidden,
		expectedBody: "{\"status\":\"error\",\"exception\":\"E-ERA-22: an authorization token is missing or wrong. " +
			"please make sure you provided a valid token.\"}",
	}
	insertRowRequest := exasol_rest_api.RowsRequest{
		SchemaName: "foo",
		TableName:  "bar",
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "C2",
				Value:      2,
			},
		},
	}
	body, err := json.Marshal(insertRowRequest)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendDeleteRows(&data, body))
}

// [itest->dsn~delete-rows-endpoint~1]
// [itest->dsn~delete-rows-request-body~1]
// [itest->dsn~delete-rows-response-body~1]
func (suite *IntegrationTestSuite) TestDeleteRowsMissingRequestParameter() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody: "{\"status\":\"error\",\"exception\":\"E-ERA-19: request has some missing parameters. " +
			"Please specify schema name, table name and condition: column name, value\"}",
	}
	request := exasol_rest_api.RowsRequest{}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendDeleteRows(&data, body))
}

// [itest->dsn~update-rows-endpoint~1]
// [itest->dsn~update-rows-request-body~1]
// [itest->dsn~update-rows-response-body~1]
func (suite *IntegrationTestSuite) TestUpdateRows() {
	username := "UPDATE_ROWS_USER"
	password := "secret"
	schemaName := "TEST_SCHEMA_UPDATE_ROWS_1"
	tableName := "ALL_DATA_TYPES"
	columns := "C1 VARCHAR(100), C2 DECIMAL(5,0), C3 BOOLEAN"

	suite.creatSchemaAndTable(schemaName, tableName, columns)
	suite.insertRowIntoTable(schemaName, tableName, "'row1', 1, true")
	suite.insertRowIntoTable(schemaName, tableName, "'row2', 2, false")
	suite.insertRowIntoTable(schemaName, tableName, "'row3', 3, false")
	suite.createExasolUser(username, password)
	suite.grantToUser(username, "CREATE SESSION")
	suite.grantToUser(username, "UPDATE ON SCHEMA "+schemaName)

	data := testData{
		server:         suite.createServerWithUser(username, password),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\"}",
	}
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: schemaName,
		TableName:  tableName,
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "C1", Value: "updated row"},
			{ColumnName: "C2", Value: 5},
		},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "C3",
				Value:      true,
			},
			ComparisonPredicate: "!=",
		},
	}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendUpdateRows(&data, body))
	suite.assertUpdatedValuesInTable(schemaName, tableName)
}

func (suite *IntegrationTestSuite) assertUpdatedValuesInTable(schemaName string, tableName string) {
	rows, err := suite.connection.Query("SELECT * FROM " + schemaName + "." + tableName)
	onError(err)
	defer func() { onError(rows.Close()) }()

	var c1 string
	var c2 int
	var c3 bool

	suite.True(rows.Next())
	err = rows.Scan(&c1, &c2, &c3)
	onError(err)
	suite.Equal("row1", c1)
	suite.Equal(1, c2)
	suite.Equal(true, c3)

	suite.True(rows.Next())
	err = rows.Scan(&c1, &c2, &c3)
	onError(err)
	suite.Equal("updated row", c1)
	suite.Equal(5, c2)
	suite.Equal(false, c3)

	suite.True(rows.Next())
	err = rows.Scan(&c1, &c2, &c3)
	onError(err)
	suite.Equal("updated row", c1)
	suite.Equal(5, c2)
	suite.Equal(false, c3)

	suite.False(rows.Next())
}

// [itest->dsn~update-rows-endpoint~1]
// [itest->dsn~update-rows-request-body~1]
// [itest->dsn~update-rows-response-body~1]
func (suite *IntegrationTestSuite) TestUpdateRowsAuthorizationError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "12345678912345678912345678912345",
		expectedStatus: http.StatusForbidden,
		expectedBody: "{\"status\":\"error\",\"exception\":\"E-ERA-22: an authorization token is missing or wrong. " +
			"please make sure you provided a valid token.\"}",
	}
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: "foo",
		TableName:  "bar",
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "C1", Value: "updated row"},
		},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "C3",
				Value:      true,
			},
		},
	}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendUpdateRows(&data, body))
}

// [itest->dsn~update-rows-endpoint~1]
// [itest->dsn~update-rows-request-body~1]
// [itest->dsn~update-rows-response-body~1]
func (suite *IntegrationTestSuite) TestUpdateRowsBadRequestError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody: "{\"status\":\"error\",\"exception\":\"E-ERA-20: update rows request has some missing parameters. " +
			"Please specify schema name, table name, values to update and condition\"}",
	}
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName:     "foo",
		TableName:      "bar",
		ValuesToUpdate: []exasol_rest_api.Value{},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "C3",
				Value:      true,
			},
		},
	}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendUpdateRows(&data, body))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRows() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE&columnName=X&value=15&valueType=int&comparisonPredicate==",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"rows\":[{\"X\":15,\"Y\":\"test\"}],\"meta\":{\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100}}]}}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsWithoutPredicate() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"rows\":[{\"X\":15,\"Y\":\"test\"},{\"X\":10,\"Y\":\"test_2\"}],\"meta\":{\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100}}]}}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsPredicateWithoutColumnName() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE&value=15&valueType=int&comparisonPredicate==",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-30: incomplete condition in the request. provide 'columnName', 'valueType' and 'value' for the condition or remove the condition\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsPredicateWithoutValue() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE&columnName=X&valueType=int&comparisonPredicate==",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-30: incomplete condition in the request. provide 'columnName', 'valueType' and 'value' for the condition or remove the condition\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsPredicateWithoutValueType() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE&columnName=X&value=10",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-30: incomplete condition in the request. provide 'columnName', 'valueType' and 'value' for the condition or remove the condition\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsWithMissingSchemaName() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "tableName=TEST_TABLE&columnName=X&value=15&valueType=int&comparisonPredicate==",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-19: request has some missing parameters. Please specify schema name and table name\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsWithMissingTableName() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=&columnName=X&value=15&valueType=int&comparisonPredicate==",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-19: request has some missing parameters. Please specify schema name and table name\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsWithIncorrectValueType() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE&columnName=X&value=15&valueType=foo&comparisonPredicate==",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-28: cannot decode value '15' with the provided value type 'foo': 'unsupported value type: foo'\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsWithNotParsableValue() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE&columnName=X&value=aaa&valueType=int&comparisonPredicate==",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-28: cannot decode value 'aaa' with the provided value type 'int': 'strconv.Atoi: parsing \\\"aaa\\\": invalid syntax'\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~get-rows-endpoint~1]
// [itest->dsn~get-rows-request-parameters~1]
// [itest->dsn~get-rows-response-body~2]
func (suite *IntegrationTestSuite) TestGetRowsWithoutAuthentication() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "schemaName=TEST_SCHEMA_1&tableName=TEST_TABLE&columnName=X&value=15&valueType=int&comparisonPredicate==",
		authToken:      "asfkndfkhjikfghsg48ghahe25nbasm32h",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-22: an authorization token is missing or wrong. please make sure you provided a valid token.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetRows(&data))
}

// [itest->dsn~execute-statement-endpoint~1]
// [itest->dsn~execute-statement-request-body~1]
// [itest->dsn~execute-statement-response-body~1]
func (suite *IntegrationTestSuite) TestExecuteStatement() {
	username := "EXECUTE_STATEMENT_USER"
	password := "secret"
	schemaName := "TEST_SCHEMA_EXECUTE_STATEMENT_1"
	columns := "C1 VARCHAR(100), C2 DECIMAL(5,0)"

	suite.creatSchemaAndTable(schemaName, "TEST_TABLE_1", columns)
	suite.createExasolUser(username, password)
	suite.grantToUser(username, "CREATE SESSION")
	suite.grantToUser(username, "CREATE ANY SCRIPT")

	data := testData{
		server:         suite.createServerWithUser(username, password),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\"}",
	}
	request := exasol_rest_api.ExecuteStatementRequest{
		Statement: "CREATE LUA SCALAR SCRIPT " + schemaName + ".hello_world () RETURNS VARCHAR(100) AS\nfunction run(ctx)\n   return 'Hello World!'\nend\n;",
	}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendExecuteStatement(&data, body))
}

// [itest->dsn~execute-statement-endpoint~1]
// [itest->dsn~execute-statement-request-body~1]
// [itest->dsn~execute-statement-response-body~1]
func (suite *IntegrationTestSuite) TestExecuteStatementAuthorizationError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "12345678912345678912345678912345",
		expectedStatus: http.StatusForbidden,
		expectedBody: "{\"status\":\"error\",\"exception\":\"E-ERA-22: an authorization token is missing or wrong. " +
			"please make sure you provided a valid token.\"}",
	}
	request := exasol_rest_api.ExecuteStatementRequest{
		Statement: "some statement",
	}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendExecuteStatement(&data, body))
}

// [itest->dsn~execute-statement-endpoint~1]
// [itest->dsn~execute-statement-request-body~1]
// [itest->dsn~execute-statement-response-body~1]
func (suite *IntegrationTestSuite) TestExecuteStatementWithSyntaxError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"error\",\"exception\":\"E-ERA-31: error while executing statement 'CREATE LUA SCALAR SCRIPT my_script;': E-EGOD-11: execution failed with SQL error code '42000' and message 'syntax error, unexpecte",
	}
	request := exasol_rest_api.ExecuteStatementRequest{
		Statement: "CREATE LUA SCALAR SCRIPT my_script;",
	}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyContains(&data, suite.sendExecuteStatement(&data, body))
}

func (suite *IntegrationTestSuite) TestRateLimiter() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "SELECT * FROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"rows\":[{\"X\":15,\"Y\":\"test\"},{\"X\":10,\"Y\":\"test_2\"}],\"meta\":{\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100}}]}}",
	}
	router := suite.startServer(data.server)

	for i := 0; i < 30; i++ {
		suite.assertResponseBodyEquals(&data, suite.sendQueryRequestWithReusableServer(&data, router))
	}

	throttledData := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "SELECT * FROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusTooManyRequests,
		expectedBody:   "Limit exceeded",
	}
	suite.assertResponseBodyEquals(&throttledData, suite.sendQueryRequestWithReusableServer(&throttledData, router))
}

func (suite *IntegrationTestSuite) sendQueryRequestWithReusableServer(data *testData,
	router *gin.Engine) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/"+data.query, nil)
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	return responseRecorder
}

type testData struct {
	query          string
	authToken      string
	expectedStatus int
	expectedBody   string
	server         exasol_rest_api.Application
}

// [itest->dsn~execute-statement-headers~1]
func (suite *IntegrationTestSuite) sendExecuteStatement(data *testData, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPost, "/api/v1/statement", bytes.NewReader(body))
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

// [itest->dsn~get-rows-headers~1]
func (suite *IntegrationTestSuite) sendGetRows(data *testData) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/rows?"+data.query, nil)
	onError(err)
	req.Header.Set("Authorization", data.authToken)
	return suite.sendHttpRequest(data, req)
}

// [itest->dsn~update-rows-headers~1]
func (suite *IntegrationTestSuite) sendUpdateRows(data *testData, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPut, "/api/v1/rows", bytes.NewReader(body))
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

// [itest->dsn~delete-rows-headers~1]
func (suite *IntegrationTestSuite) sendDeleteRows(data *testData, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodDelete, "/api/v1/rows", bytes.NewReader(body))
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

// [itest->dsn~insert-row-headers~1]
func (suite *IntegrationTestSuite) sendInsertRow(data *testData, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPost, "/api/v1/row", bytes.NewReader(body))
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

// [itest->dsn~get-tables-headers~1]
func (suite *IntegrationTestSuite) sendGetTables(data *testData) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/tables", nil)
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

// [itest->dsn~execute-query-request-parameters~1]
// [itest->dsn~execute-query-headers~1]
func (suite *IntegrationTestSuite) sendQueryRequest(data *testData) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/"+data.query, nil)
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

func (suite *IntegrationTestSuite) sendHttpRequest(data *testData, req *http.Request) *httptest.ResponseRecorder {
	router := suite.startServer(data.server)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	return responseRecorder
}

func (suite *IntegrationTestSuite) assertResponseBodyEquals(data *testData, responseRecorder *httptest.ResponseRecorder) {
	suite.Equal(data.expectedStatus, responseRecorder.Code)
	suite.Equal(data.expectedBody, responseRecorder.Body.String())
}

func (suite *IntegrationTestSuite) assertResponseBodyContains(data *testData, responseRecorder *httptest.ResponseRecorder) {
	suite.Equal(data.expectedStatus, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(), data.expectedBody)
}

func (suite *IntegrationTestSuite) assertTableHasOnlyOneRow(schemaName string, tableName string) {
	rows, err := suite.connection.Query("SELECT * FROM " + schemaName + "." + tableName)
	onError(err)
	defer func() { onError(rows.Close()) }()
	suite.True(rows.Next())
	suite.False(rows.Next())
}

func runExasolContainer() *testSetupAbstraction.TestSetupAbstraction {
	dbVersion := os.Getenv("EXASOL_DB_VERSION")
	if dbVersion == "" {
		dbVersion = "8.34.0"
	}
	exasolContainer, err := testSetupAbstraction.New().CloudSetupConfigFilePath("no-config.json").DockerDbVersion(dbVersion).Start()
	onError(err)
	return exasolContainer
}

func onError(err error) {
	if err != nil {
		log.Printf("Error %s", err)
		panic(err)
	}
}

func (suite *IntegrationTestSuite) createServerWithUser(user string, password string) exasol_rest_api.Application {
	properties := &exasol_rest_api.ApplicationProperties{
		APITokens:                       suite.defaultAuthTokens,
		ExasolUser:                      user,
		ExasolPassword:                  password,
		ExasolHost:                      suite.exasolHost,
		ExasolPort:                      suite.exasolPort,
		ExasolValidateServerCertificate: "false",
	}
	return suite.runApiServer(properties)
}
func (suite *IntegrationTestSuite) createServerWithDefaultProperties() exasol_rest_api.Application {
	return suite.createServerWithUser(suite.defaultServiceUsername, suite.defaultServicePassword)
}

func (suite *IntegrationTestSuite) runApiServer(properties *exasol_rest_api.ApplicationProperties) exasol_rest_api.Application {
	return exasol_rest_api.Application{
		Properties: properties,
		Authorizer: &exasol_rest_api.TokenAuthorizer{
			AllowedTokens: exasol_rest_api.CreateStringsSet(properties.APITokens),
		},
	}
}

func createDefaultServiceUserWithAccess(database *sql.DB, user, password string) {
	schemaName := "TEST_SCHEMA_1"
	_, err := database.Exec("CREATE SCHEMA " + schemaName)
	onError(err)
	_, err = database.Exec("CREATE TABLE " + schemaName + ".TEST_TABLE(X INT, Y VARCHAR(100))")
	onError(err)
	_, err = database.Exec("INSERT INTO " + schemaName + ".TEST_TABLE VALUES (15, 'test')")
	onError(err)
	_, err = database.Exec("INSERT INTO " + schemaName + ".TEST_TABLE VALUES (10, 'test_2')")
	onError(err)

	_, err = database.Exec("CREATE USER " + user + " IDENTIFIED BY \"" + password + "\"")
	onError(err)
	_, err = database.Exec("GRANT CREATE SESSION TO " + user)
	onError(err)
	_, err = database.Exec("GRANT SELECT ON SCHEMA " + schemaName + " TO " + user)
	onError(err)
	_, err = database.Exec("COMMIT")
	onError(err)
}

func (suite *IntegrationTestSuite) createExasolUser(username string, password string) {
	_, err := suite.connection.Exec("CREATE USER " + username + " IDENTIFIED BY \"" + password + "\"")
	onError(err)
}

func (suite *IntegrationTestSuite) grantToUser(username string, privilege string) {
	_, err := suite.connection.Exec("GRANT " + privilege + " TO " + username)
	onError(err)
}

func (suite *IntegrationTestSuite) creatSchemaAndTable(schemaName string, tableName string, columns string) {
	_, err := suite.connection.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName)
	onError(err)
	_, err = suite.connection.Exec("CREATE TABLE " + schemaName + "." + tableName + "(" + columns + ")")
	onError(err)
}

func (suite *IntegrationTestSuite) insertRowIntoTable(schemaName string, tableName string, values string) {
	_, err := suite.connection.Exec("INSERT INTO " + schemaName + "." + tableName + " VALUES (" + values + ")")
	onError(err)
}
