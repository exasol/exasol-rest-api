package exasol_rest_api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/exasol/exasol-driver-go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type IntegrationTestSuite struct {
	suite.Suite
	ctx                    context.Context
	exasolContainer        testcontainers.Container
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
	suite.ctx = context.Background()
	suite.defaultServiceUsername = "api_service_account"
	suite.defaultServicePassword = "secret_password"
	suite.defaultAuthTokens = []string{"3J90XAv9loMIXzQdfYmtJrHAbopPsc", "OR6rq6KjWmhvGU770A9OTjpfH86nlk"}
	suite.exasolContainer = runExasolContainer(suite.ctx)
	suite.exasolHost = getExasolHost(suite.exasolContainer, suite.ctx)
	suite.exasolPort = 8563
	database, err := sql.Open("exasol",
		exasol.NewConfig("sys",
			"exasol").UseTLS(false).Host(suite.exasolHost).Port(suite.exasolPort).Autocommit(true).String())
	onError(err)
	suite.connection = database
	createDefaultServiceUserWithAccess(suite.defaultServiceUsername, suite.defaultServicePassword, suite.exasolHost,
		suite.exasolPort)
}

func (suite *IntegrationTestSuite) startServer(application exasol_rest_api.Application) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	exasol_rest_api.AddEndpoints(router, application)
	suite.appProperties = application.Properties
	return router
}

func (suite *IntegrationTestSuite) TestQuery() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "SELECT * FROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":1,\"numRowsInMessage\":1,\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100,\"characterSet\":\"UTF8\"}}],\"data\":[[15],[\"test\"]]}}],\"numResults\":1}}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestQueryWithTypo() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "SELECTFROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"error\",\"exception\":{\"text\":\"syntax error, unexpected ",
	}
	suite.assertResponseBodyContains(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestInsertNotAllowed() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "CREATE SCHEMA not_allowed_schema",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"error\",\"exception\":{\"text\":\"insufficient privileges for creating schema",
	}
	suite.assertResponseBodyContains(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestExasolUserWithoutCreateSessionPrivilege() {
	username := "user_without_session_privilege"
	password := "secret"
	suite.createExasolUser(username, password)

	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                username,
		ExasolPassword:            password,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})
	data := testData{
		server:         server,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol: [08004] Connection exception - insufficient privileges: CREATE SESSION.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestExasolUserWithWrongCredentials() {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                "not_existing_user",
		ExasolPassword:            "wrong_password",
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})
	data := testData{
		server:         server,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol: [08004] Connection exception - authentication failed.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestWrongExasolPort() {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                suite.defaultServiceUsername,
		ExasolPassword:            suite.defaultServicePassword,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                4321,
		ExasolWebsocketAPIVersion: 2,
	})
	data := testData{
		server:         server,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol:",
	}
	suite.assertResponseBodyContains(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestWrongWebsocketApiVersion() {
	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                suite.defaultServiceUsername,
		ExasolPassword:            suite.defaultServicePassword,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 0,
	})
	data := testData{
		server:         server,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol: E-ERA-15: error while sending a login command via websockets connection: [00000] Could not create WebSocket protocol version 0\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestUnauthorizedAccessToQuery() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "some query",
		authToken:      "OR6rq6KjWmhvGU770A9OTjpfH86nlkq",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"Error\":\"E-ERA-22: an authorization token is missing or wrong. please make sure you provided a valid token.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestUnauthorizedAccessWithShortToken() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		query:          "some query",
		authToken:      "tooshort",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"Error\":\"E-ERA-23: an authorization token has invalid length: 8. please only use tokens with the length longer or equal to 30.\"}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendQueryRequest(&data))
}

func (suite *IntegrationTestSuite) TestGetTables() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":10,\"numRows\":0,\"numRowsInMessage\":0,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_OWNER\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_OBJECT_ID\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"TABLE_IS_VIRTUAL\",\"dataType\":{\"type\":\"BOOLEAN\"}},{\"name\":\"TABLE_HAS_DISTRIBUTION_KEY\",\"dataType\":{\"type\":\"BOOLEAN\"}},{\"name\":\"TABLE_HAS_PARTITION_KEY\",\"dataType\":{\"type\":\"BOOLEAN\"}},{\"name\":\"TABLE_ROW_COUNT\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"DELETE_PERCENTAGE\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":4,\"scale\":1}},{\"name\":\"TABLE_COMMENT\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":2000,\"characterSet\":\"UTF8\"}}]}}],\"numResults\":1}}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetTables(&data))
}

func (suite *IntegrationTestSuite) TestGetTablesWithZeroTables() {
	username := "USER_WITHOUT_OWNED_SCHEMA"
	password := "secret"
	suite.createExasolUser(username, password)
	suite.grantToUser(username, "CREATE SESSION")

	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                username,
		ExasolPassword:            password,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})
	data := testData{
		server:         server,
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":10,\"numRows\":0,\"numRowsInMessage\":0,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_OWNER\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_OBJECT_ID\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"TABLE_IS_VIRTUAL\",\"dataType\":{\"type\":\"BOOLEAN\"}},{\"name\":\"TABLE_HAS_DISTRIBUTION_KEY\",\"dataType\":{\"type\":\"BOOLEAN\"}},{\"name\":\"TABLE_HAS_PARTITION_KEY\",\"dataType\":{\"type\":\"BOOLEAN\"}},{\"name\":\"TABLE_ROW_COUNT\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"DELETE_PERCENTAGE\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":4,\"scale\":1}},{\"name\":\"TABLE_COMMENT\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":2000,\"characterSet\":\"UTF8\"}}]}}],\"numResults\":1}}",
	}
	suite.assertResponseBodyEquals(&data, suite.sendGetTables(&data))
}

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

	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                username,
		ExasolPassword:            password,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})

	data := testData{
		server:         server,
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"rowCount\",\"rowCount\":1}],\"numResults\":1}}",
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

func (suite *IntegrationTestSuite) assertInsertRowValuesInTable(schemaName string, tableName string) {
	rows, err := suite.connection.Query("SELECT * FROM " + schemaName + "." + tableName)
	onError(err)
	defer rows.Close()
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

func (suite *IntegrationTestSuite) TestInsertRowAuthorizationError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "badToken",
		expectedStatus: http.StatusForbidden,
		expectedBody: "{\"Error\":\"E-ERA-23: an authorization token has invalid length: 8. " +
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

func (suite *IntegrationTestSuite) TestInsertRowMissingRequestParameter() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody: "{\"Error\":\"E-ERA-17: insert row request has some missing parameters. " +
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

	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                username,
		ExasolPassword:            password,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})

	data := testData{
		server:         server,
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"rowCount\",\"rowCount\":2}],\"numResults\":1}}",
	}
	deleteRowsRequest := exasol_rest_api.DeleteRowsRequest{
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

func (suite *IntegrationTestSuite) TestDeleteRowsAuthorizationError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "12345678912345678912345678912345",
		expectedStatus: http.StatusForbidden,
		expectedBody: "{\"Error\":\"E-ERA-22: an authorization token is missing or wrong. " +
			"please make sure you provided a valid token.\"}",
	}
	insertRowRequest := exasol_rest_api.DeleteRowsRequest{
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

func (suite *IntegrationTestSuite) TestDeleteRowsMissingRequestParameter() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody: "{\"Error\":\"E-ERA-19: delete rows request has some missing parameters. " +
			"Please specify schema name, table name and condition: column name, value\"}",
	}
	request := exasol_rest_api.DeleteRowsRequest{}
	body, err := json.Marshal(request)
	onError(err)
	suite.assertResponseBodyEquals(&data, suite.sendDeleteRows(&data, body))
}

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

	server := suite.runApiServer(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                username,
		ExasolPassword:            password,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})

	data := testData{
		server:         server,
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"rowCount\",\"rowCount\":2}],\"numResults\":1}}",
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
	defer rows.Close()

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

func (suite *IntegrationTestSuite) TestUpdateRowsAuthorizationError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "12345678912345678912345678912345",
		expectedStatus: http.StatusForbidden,
		expectedBody: "{\"Error\":\"E-ERA-22: an authorization token is missing or wrong. " +
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

func (suite *IntegrationTestSuite) TestUpdateRowsBadRequestError() {
	data := testData{
		server:         suite.createServerWithDefaultProperties(),
		authToken:      "foo",
		expectedStatus: http.StatusBadRequest,
		expectedBody: "{\"Error\":\"E-ERA-20: update rows request has some missing parameters. " +
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

type testData struct {
	query          string
	authToken      string
	expectedStatus int
	expectedBody   string
	server         exasol_rest_api.Application
}

func (suite *IntegrationTestSuite) sendUpdateRows(data *testData, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPut, "/api/v1/rows", bytes.NewReader(body))
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

func (suite *IntegrationTestSuite) sendDeleteRows(data *testData, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodDelete, "/api/v1/rows", bytes.NewReader(body))
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

func (suite *IntegrationTestSuite) sendInsertRow(data *testData, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPost, "/api/v1/row", bytes.NewReader(body))
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

func (suite *IntegrationTestSuite) sendGetTables(data *testData) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/tables", nil)
	req.Header.Set("Authorization", data.authToken)
	onError(err)
	return suite.sendHttpRequest(data, req)
}

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

func (suite *IntegrationTestSuite) assertResponseBodyEquals(data *testData,
	responseRecorder *httptest.ResponseRecorder) {
	suite.Equal(data.expectedStatus, responseRecorder.Code)
	suite.Equal(data.expectedBody, responseRecorder.Body.String())
}

func (suite *IntegrationTestSuite) assertResponseBodyContains(data *testData,
	responseRecorder *httptest.ResponseRecorder) {
	suite.Equal(data.expectedStatus, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(), data.expectedBody)
}

func (suite *IntegrationTestSuite) assertTableHasOnlyOneRow(schemaName string, tableName string) {
	rows, err := suite.connection.Query("SELECT * FROM " + schemaName + "." + tableName)
	onError(err)
	defer rows.Close()
	suite.True(rows.Next())
	suite.False(rows.Next())
}

func runExasolContainer(ctx context.Context) testcontainers.Container {
	dbVersion := os.Getenv("DB_VERSION")
	if dbVersion == "" {
		dbVersion = "7.1.1"
	}
	request := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("exasol/docker-db:%s", dbVersion),
		ExposedPorts: []string{"8563", "2580"},
		WaitingFor:   wait.ForLog("All stages finished").WithStartupTimeout(time.Minute * 5),
		Privileged:   true,
	}
	exasolContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	onError(err)
	return exasolContainer
}

func getExasolHost(exasolContainer testcontainers.Container, ctx context.Context) string {
	host, err := exasolContainer.ContainerIP(ctx)
	onError(err)
	return host
}

func onError(err error) {
	if err != nil {
		log.Printf("Error %s", err)
		panic(err)
	}
}

func (suite *IntegrationTestSuite) createServerWithDefaultProperties() exasol_rest_api.Application {
	properties := &exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                suite.defaultServiceUsername,
		ExasolPassword:            suite.defaultServicePassword,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		Encryption:                false,
		UseTLS:                    false,
		ExasolWebsocketAPIVersion: 2,
	}
	return suite.runApiServer(properties)
}

func (suite *IntegrationTestSuite) runApiServer(properties *exasol_rest_api.ApplicationProperties) exasol_rest_api.Application {
	return exasol_rest_api.Application{
		Properties: properties,
		Authorizer: &exasol_rest_api.TokenAuthorizer{
			AllowedTokens: exasol_rest_api.CreateStringsSet(properties.APITokens),
		},
	}
}

func createDefaultServiceUserWithAccess(user string, password string, host string, port int) {
	database, err := sql.Open("exasol",
		exasol.NewConfig("sys", "exasol").UseTLS(false).Host(host).Port(port).Autocommit(true).String())
	onError(err)
	schemaName := "TEST_SCHEMA_1"
	_, err = database.Exec("CREATE SCHEMA " + schemaName)
	onError(err)
	_, err = database.Exec("CREATE TABLE " + schemaName + ".TEST_TABLE(x INT, y VARCHAR(100))")
	onError(err)
	_, err = database.Exec("INSERT INTO " + schemaName + ".TEST_TABLE VALUES (15, 'test')")
	onError(err)

	_, err = database.Exec("CREATE USER " + user + " IDENTIFIED BY \"" + password + "\"")
	onError(err)
	_, err = database.Exec("GRANT CREATE SESSION TO " + user)
	onError(err)
	_, err = database.Exec("GRANT SELECT ON SCHEMA " + schemaName + " TO " + user)
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
	_, err := suite.connection.Exec("CREATE SCHEMA " + schemaName)
	onError(err)
	_, err = suite.connection.Exec("CREATE TABLE " + schemaName + "." + tableName + "(" + columns + ")")
	onError(err)
}

func (suite *IntegrationTestSuite) insertRowIntoTable(schemaName string, tableName string, values string) {
	_, err := suite.connection.Exec("INSERT INTO " + schemaName + "." + tableName + " VALUES (" + values + ")")
	onError(err)
}
