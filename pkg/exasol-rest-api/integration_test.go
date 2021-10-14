package exasol_rest_api_test

import (
	"context"
	"database/sql"
	"log"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"net/http"
	"net/http/httptest"
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
	ctx                   context.Context
	exasolContainer       testcontainers.Container
	defaultExasolUsername string
	defaultExasolPassword string
	defaultAuthTokens     []string
	exasolPort            int
	exasolHost            string
	appProperties         *exasol_rest_api.ApplicationProperties
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.defaultExasolUsername = "api_service_account"
	suite.defaultExasolPassword = "secret_password"
	suite.defaultAuthTokens = []string{"3J90XAv9loMIXzQdfYmtJrHAbopPsc", "OR6rq6KjWmhvGU770A9OTjpfH86nlk"}
	suite.exasolContainer = runExasolContainer(suite.ctx)
	suite.exasolHost = getExasolHost(suite.exasolContainer, suite.ctx)
	suite.exasolPort = 8563
	createDefaultServiceUserWithAccess(suite.defaultExasolUsername, suite.defaultExasolPassword, suite.exasolHost,
		suite.exasolPort)
}

func (suite *IntegrationTestSuite) startServer(application exasol_rest_api.Application) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/query/:query", application.Query)
	suite.appProperties = application.Properties
	return router
}

func (suite *IntegrationTestSuite) TestQuery() {
	data := testData{
		application:    suite.createApplicationWithDefaultProperties(),
		query:          "SELECT * FROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":1,\"numRowsInMessage\":1,\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100,\"characterSet\":\"UTF8\"}}],\"data\":[[15],[\"test\"]]}}],\"numResults\":1}}",
	}
	suite.validateResponseBodyEquals(&data)
}

func (suite *IntegrationTestSuite) TestQueryWithTypo() {
	data := testData{
		application:    suite.createApplicationWithDefaultProperties(),
		query:          "SELECTFROM TEST_SCHEMA_1.TEST_TABLE",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"error\",\"exception\":{\"text\":\"syntax error, unexpected IDENTIFIER_LIST_",
	}
	suite.validateResponseBodyContains(&data)
}

func (suite *IntegrationTestSuite) TestInsertNotAllowed() {
	data := testData{
		application:    suite.createApplicationWithDefaultProperties(),
		query:          "CREATE SCHEMA not_allowed_schema",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusOK,
		expectedBody:   "{\"status\":\"error\",\"exception\":{\"text\":\"insufficient privileges for creating schema",
	}
	suite.validateResponseBodyContains(&data)
}

func (suite *IntegrationTestSuite) TestExasolUserWithoutCreateSessionPrivilege() {
	username := "user_without_session_privilege"
	password := "secret"
	suite.createExasolUser(username, password)

	application := suite.createApplication(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                username,
		ExasolPassword:            password,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})
	data := testData{
		application:    application,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol: [08004] Connection exception - insufficient privileges: CREATE SESSION.\"}",
	}
	suite.validateResponseBodyEquals(&data)
}

func (suite *IntegrationTestSuite) TestExasolUserWithWrongCredentials() {
	application := suite.createApplication(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                "not_existing_user",
		ExasolPassword:            "wrong_password",
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 2,
	})
	data := testData{
		application:    application,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol: [08004] Connection exception - authentication failed.\"}",
	}
	suite.validateResponseBodyEquals(&data)
}

func (suite *IntegrationTestSuite) TestWrongExasolPort() {
	application := suite.createApplication(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                suite.defaultExasolUsername,
		ExasolPassword:            suite.defaultExasolPassword,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                4321,
		ExasolWebsocketAPIVersion: 2,
	})
	data := testData{
		application:    application,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol:",
	}
	suite.validateResponseBodyContains(&data)
}

func (suite *IntegrationTestSuite) TestWrongWebsocketApiVersion() {
	application := suite.createApplication(&exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                suite.defaultExasolUsername,
		ExasolPassword:            suite.defaultExasolPassword,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketAPIVersion: 0,
	})
	data := testData{
		application:    application,
		query:          "some query",
		authToken:      suite.defaultAuthTokens[0],
		expectedStatus: http.StatusBadRequest,
		expectedBody:   "{\"Error\":\"E-ERA-2: error while opening a connection with Exasol: E-ERA-15: error while sending a login command via websockets connection: [00000] Could not create WebSocket protocol version 0\"}",
	}
	suite.validateResponseBodyEquals(&data)
}

func (suite *IntegrationTestSuite) TestUnauthorizedAccessToQuery() {
	data := testData{
		application:    suite.createApplicationWithDefaultProperties(),
		query:          "some query",
		authToken:      "OR6rq6KjWmhvGU770A9OTjpfH86nlkq",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"Error\":\"E-ERA-22: an authorization token is missing or wrong. please make sure you provided a valid token.\"}",
	}
	suite.validateResponseBodyEquals(&data)
}

func (suite *IntegrationTestSuite) TestUnauthorizedAccessWithShortToken() {
	data := testData{
		application:    suite.createApplicationWithDefaultProperties(),
		query:          "some query",
		authToken:      "tooshort",
		expectedStatus: http.StatusForbidden,
		expectedBody:   "{\"Error\":\"E-ERA-23: an authorization token has invalid length: 8. please only use tokens with the length longer or equal to 30.\"}",
	}
	suite.validateResponseBodyEquals(&data)
}

type testData struct {
	query          string
	authToken      string
	expectedStatus int
	expectedBody   string
	application    exasol_rest_api.Application
}

func (suite *IntegrationTestSuite) validateResponseBodyEquals(data *testData) {
	router := suite.startServer(data.application)
	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/"+data.query, nil)
	req.Header.Set("Authorization", data.authToken)
	onError(err)

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(data.expectedStatus, responseRecorder.Code)
	suite.Equal(data.expectedBody, responseRecorder.Body.String())
}

func (suite *IntegrationTestSuite) validateResponseBodyContains(data *testData) {
	router := suite.startServer(data.application)
	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/"+data.query, nil)
	req.Header.Set("Authorization", data.authToken)
	onError(err)

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(data.expectedStatus, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(), data.expectedBody)
}

func runExasolContainer(ctx context.Context) testcontainers.Container {
	request := testcontainers.ContainerRequest{
		Image:        "exasol/docker-db:7.1.1",
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

func (suite *IntegrationTestSuite) createApplicationWithDefaultProperties() exasol_rest_api.Application {
	properties := &exasol_rest_api.ApplicationProperties{
		APITokens:                 suite.defaultAuthTokens,
		ExasolUser:                suite.defaultExasolUsername,
		ExasolPassword:            suite.defaultExasolPassword,
		ExasolHost:                suite.exasolHost,
		ExasolPort:                suite.exasolPort,
		Encryption:                false,
		UseTLS:                    false,
		ExasolWebsocketAPIVersion: 2,
	}
	return suite.createApplication(properties)
}

func (suite *IntegrationTestSuite) createApplication(properties *exasol_rest_api.ApplicationProperties) exasol_rest_api.Application {
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
	database, err := sql.Open("exasol", exasol.NewConfig("sys", "exasol").UseTLS(false).
		Host(suite.exasolHost).Port(suite.exasolPort).String())
	onError(err)
	_, err = database.Exec("CREATE USER " + username + " IDENTIFIED BY \"" + password + "\"")
	onError(err)
}
