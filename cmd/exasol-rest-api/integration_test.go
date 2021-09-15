package exasol_rest_api_test

import (
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	exasol_rest_api "main/cmd/exasol-rest-api"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
}

func (suite *IntegrationTestSuite) TestGetMethod() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/query/:query", exasol_rest_api.Query)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/query/SELECT 1 FROM DUAL", nil)
	onError(err)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	suite.Equal(http.StatusOK, responseRecorder.Code)
	suite.Equal("\"{\\\"status\\\":\\\"ok\\\",\\\"responseData\\\":{\\\"results\\\":[{\\\"resultType\\\":\\\"resultSet\\\",\\\"resultSet\\\":{\\\"numColumns\\\":1,\\\"numRows\\\":1,\\\"numRowsInMessage\\\":1,\\\"columns\\\":[{\\\"name\\\":\\\"1\\\",\\\"dataType\\\":{\\\"type\\\":\\\"DECIMAL\\\",\\\"precision\\\":1,\\\"scale\\\":0}}],\\\"data\\\":[[1]]}}],\\\"numResults\\\":1}}\"",
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
