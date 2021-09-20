package exasol_rest_api_test

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"net/http"
	"net/http/httptest"
	"testing"

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
	exasolPort            int
	appProperties         *exasol_rest_api.ApplicationProperties
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.defaultExasolUsername = "api_service_account"
	suite.defaultExasolPassword = "secret_password"
	suite.exasolContainer = runExasolContainer(suite.ctx)
	suite.exasolPort = getExasolPort(suite.exasolContainer, suite.ctx)
	createDefaultServiceUserWithAccess(suite.defaultExasolUsername, suite.defaultExasolPassword, suite.exasolPort)
}

func (suite *IntegrationTestSuite) startServer(application exasol_rest_api.Application) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/query/:query", application.Query)
	suite.appProperties = application.Properties
	return router
}

func (suite *IntegrationTestSuite) TestQuery() {
	router := suite.startServer(suite.createApplicationWithDefaultProperties())
	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusOK, responseRecorder.Code)
	suite.Equal("{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":1,\"numRowsInMessage\":1,\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100,\"characterSet\":\"UTF8\"}}],\"data\":[[15],[\"test\"]]}}],\"numResults\":1}}",
		string(responseRecorder.Body.Bytes()))
}

func (suite *IntegrationTestSuite) TestInsertNotAllowed() {
	router := suite.startServer(suite.createApplicationWithDefaultProperties())
	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/CREATE SCHEMA not_allowed_schema", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusOK, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(),
		"{\"status\":\"error\",\"exception\":{\"text\":\"insufficient privileges for creating schema")
}

func (suite *IntegrationTestSuite) TestExasolUserWithoutCreateSessionPrivilege() {
	username := "user_without_session"
	password := "secret"
	suite.createExasolUser(username, password)
	router := suite.startServer(suite.createApplication(&exasol_rest_api.ApplicationProperties{
		ExasolUser:                username,
		ExasolPassword:            password,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketApiVersion: 2,
	}))

	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusBadRequest, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(),
		"{\"ErrorCode\":\"EXA-REST-API-1\",\"Message\":\"[08004] Connection exception - insufficient privileges: CREATE SESSION.\"}")
}

func (suite *IntegrationTestSuite) TestExasolUserWithWrongCredentials() {
	router := suite.startServer(suite.createApplication(&exasol_rest_api.ApplicationProperties{
		ExasolUser:                "not_existing_user",
		ExasolPassword:            "wrong_password",
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketApiVersion: 2,
	}))

	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusBadRequest, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(),
		"{\"ErrorCode\":\"EXA-REST-API-1\",\"Message\":\"[08004] Connection exception - authentication failed.\"}")
}

func (suite *IntegrationTestSuite) TestWrongExasolPort() {
	router := suite.startServer(suite.createApplication(&exasol_rest_api.ApplicationProperties{
		ExasolUser:                suite.defaultExasolUsername,
		ExasolPassword:            suite.defaultExasolPassword,
		ExasolPort:                4321,
		ExasolWebsocketApiVersion: 2,
	}))

	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusBadRequest, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(), "{\"ErrorCode\":\"EXA-REST-API-1\"")
	suite.Contains(responseRecorder.Body.String(), "connect: connection refused")
}

func (suite *IntegrationTestSuite) TestWrongWebsocketApiVersion() {
	router := suite.startServer(suite.createApplication(&exasol_rest_api.ApplicationProperties{
		ExasolUser:                suite.defaultExasolUsername,
		ExasolPassword:            suite.defaultExasolPassword,
		ExasolPort:                suite.exasolPort,
		ExasolWebsocketApiVersion: 0,
	}))

	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusBadRequest, responseRecorder.Code)
	suite.Contains(responseRecorder.Body.String(),
		"{\"ErrorCode\":\"EXA-REST-API-1\",\"Message\":\"[00000] Could not create WebSocket protocol version 0\"}")
}

func runExasolContainer(ctx context.Context) testcontainers.Container {
	request := testcontainers.ContainerRequest{
		Image:        "exasol/docker-db:7.0.10",
		ExposedPorts: []string{"8563", "2580"},
		WaitingFor:   wait.ForLog("All stages finished"),
		Privileged:   true,
	}
	exasolContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	onError(err)
	return exasolContainer
}

func getExasolPort(exasolContainer testcontainers.Container, ctx context.Context) int {
	port, err := exasolContainer.MappedPort(ctx, "8563")
	onError(err)
	return port.Int()
}

func onError(err error) {
	if err != nil {
		log.Printf("Error %s", err)
		panic(err)
	}
}

func (suite *IntegrationTestSuite) createApplicationWithDefaultProperties() exasol_rest_api.Application {
	properties := &exasol_rest_api.ApplicationProperties{
		ExasolUser:                suite.defaultExasolUsername,
		ExasolPassword:            suite.defaultExasolPassword,
		ExasolHost:                "localhost",
		ExasolPort:                suite.exasolPort,
		Encryption:                false,
		UseTLS:                    false,
		ExasolWebsocketApiVersion: 2,
	}
	return exasol_rest_api.Application{
		Properties: properties,
	}
}

func (suite *IntegrationTestSuite) createApplication(properties *exasol_rest_api.ApplicationProperties) exasol_rest_api.Application {
	return exasol_rest_api.Application{
		Properties: properties,
	}
}

func createDefaultServiceUserWithAccess(user string, password string, port int) {
	database, _ := sql.Open("exasol", exasol.NewConfig("sys", "exasol").UseTLS(false).Port(port).String())
	schemaName := "TEST_SCHEMA_1"
	_, _ = database.Exec("CREATE SCHEMA " + schemaName)
	_, _ = database.Exec("CREATE TABLE " + schemaName + ".TEST_TABLE(x INT, y VARCHAR(100))")
	_, _ = database.Exec("INSERT INTO " + schemaName + ".TEST_TABLE VALUES (15, 'test')")

	_, _ = database.Exec("CREATE USER " + user + " IDENTIFIED BY \"" + password + "\"")
	_, _ = database.Exec("GRANT CREATE SESSION TO " + user)
	_, _ = database.Exec("GRANT SELECT ON SCHEMA " + schemaName + " TO " + user)
}

func (suite *IntegrationTestSuite) createExasolUser(username string, password string) {
	database, _ := sql.Open("exasol", exasol.NewConfig("sys", "exasol").UseTLS(false).Port(suite.exasolPort).String())
	_, _ = database.Exec("CREATE USER " + username + " IDENTIFIED BY \"" + password + "\"")
}
