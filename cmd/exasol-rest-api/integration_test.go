package exasol_rest_api_test

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	exasol_rest_api "main/cmd/exasol-rest-api"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/exasol/exasol-driver-go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type IntegrationTestSuite struct {
	suite.Suite
	ctx             context.Context
	exasolContainer testcontainers.Container
	port            int
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.exasolContainer = runExasolContainer(suite.ctx)
	suite.port = getExasolPort(suite.exasolContainer, suite.ctx)
	createConnectionPropertiesFile(suite)
	createTableInExasol(suite)
}

func (suite *IntegrationTestSuite) TestGetMethod() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/query/:query", exasol_rest_api.Query)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusOK, responseRecorder.Code)
	suite.Equal("\"{\\\"status\\\":\\\"ok\\\",\\\"responseData\\\":{\\\"results\\\":[{\\\"resultType\\\":\\\"resultSet\\\",\\\"resultSet\\\":{\\\"numColumns\\\":2,\\\"numRows\\\":1,\\\"numRowsInMessage\\\":1,\\\"columns\\\":[{\\\"name\\\":\\\"X\\\",\\\"dataType\\\":{\\\"type\\\":\\\"DECIMAL\\\",\\\"precision\\\":18,\\\"scale\\\":0}},{\\\"name\\\":\\\"Y\\\",\\\"dataType\\\":{\\\"type\\\":\\\"VARCHAR\\\",\\\"size\\\":100,\\\"characterSet\\\":\\\"UTF8\\\"}}],\\\"data\\\":[[15],[\\\"test\\\"]]}}],\\\"numResults\\\":1}}\"",
		responseRecorder.Body.String())
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

func createConnectionPropertiesFile(suite *IntegrationTestSuite) {
	connProperties := &exasol_rest_api.ConnectionProperties{
		User:       "sys",
		Password:   "exasol",
		Host:       "localhost",
		Port:       suite.port,
		Encryption: false,
		UseTLS:     false,
		ApiVersion: 2,
	}
	file, err := ioutil.TempFile("", "connection_properties_*.yml")
	onError(err)
	data, err := yaml.Marshal(&connProperties)
	onError(err)
	_, err = file.Write(data)
	onError(err)
	err = os.Setenv("CONNECTION_PROPERTIES_PATH", file.Name())
	onError(err)
}

func createTableInExasol(suite *IntegrationTestSuite) {
	database, _ := sql.Open("exasol", exasol.NewConfig("sys", "exasol").UseTLS(false).Port(suite.port).String())
	schemaName := "TEST_SCHEMA_1"
	_, _ = database.Exec("CREATE SCHEMA " + schemaName)
	_, _ = database.Exec("CREATE TABLE " + schemaName + ".TEST_TABLE(x INT, y VARCHAR(100))")
	_, _ = database.Exec("INSERT INTO " + schemaName + ".TEST_TABLE VALUES (15, 'test')")
}
