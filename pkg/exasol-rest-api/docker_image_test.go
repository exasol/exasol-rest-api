package exasol_rest_api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	exasol_rest_api "github.com/exasol/exasol-rest-api/pkg/exasol-rest-api"

	testSetupAbstraction "github.com/exasol/exasol-test-setup-abstraction-server/go-client"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/yaml.v3"
)

type DockerImageTestSuite struct {
	suite.Suite
	ctx                   context.Context
	exasolContainer       *testSetupAbstraction.TestSetupAbstraction
	defaultExasolUsername string
	defaultExasolPassword string
	defaultAuthTokens     string
	exasolPort            int
	exasolHost            string
}

func TestDockerImageSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(DockerImageTestSuite))
}

func (suite *DockerImageTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.defaultExasolUsername = "api_service_account"
	suite.defaultExasolPassword = "secret_password"
	suite.defaultAuthTokens = "3J90XAv9loMIXzQdfYmtJrHAbopPsc,OR6rq6KjWmhvGU770A9OTjpfH86nlk"
	suite.exasolContainer = runExasolContainer()
	connectionInfo, err := suite.exasolContainer.GetConnectionInfo()
	onError(err)
	suite.exasolHost = connectionInfo.Host
	suite.exasolPort = connectionInfo.Port
	connection, err := suite.exasolContainer.CreateConnection()
	onError(err)
	createDefaultServiceUserWithAccess(connection, suite.defaultExasolUsername, suite.defaultExasolPassword)
}

func (suite *DockerImageTestSuite) TestQueryDocker() {
	apiContainer := runRestAPIContainer(suite.restAPIProperties(), suite.exasolPort)
	ip, err := apiContainer.Host(suite.ctx)
	onError(err)

	port, err := apiContainer.MappedPort(suite.ctx, "8080")
	onError(err)

	baseUrl := "http://" + ip + ":" + port.Port()
	req, err := http.NewRequest(http.MethodGet, baseUrl+"/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	req.Header.Set("Authorization", "3J90XAv9loMIXzQdfYmtJrHAbopPsc")
	onError(err)

	client := http.Client{}
	response, err := client.Do(req)
	onError(err)

	body, err := io.ReadAll(response.Body)
	onError(err)

	suite.Equal("200 OK", response.Status)
	suite.Equal("{\"status\":\"ok\",\"rows\":[{\"X\":15,\"Y\":\"test\"},{\"X\":10,\"Y\":\"test_2\"}],\"meta\":{\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100}}]}}",
		string(body))
}

func (suite *DockerImageTestSuite) TestSwaggerDocumentationContainsProjectVersion() {
	apiContainer := runRestAPIContainer(suite.restAPIProperties(), suite.exasolPort)
	ip, err := apiContainer.Host(suite.ctx)
	suite.Require().NoError(err)
	port, err := apiContainer.MappedPort(suite.ctx, "8080")
	suite.Require().NoError(err)

	response, err := http.Get("http://" + ip + ":" + port.Port() + "/swagger/doc.json")
	suite.Require().NoError(err)
	defer func() { suite.NoError(response.Body.Close()) }()
	suite.Equal(http.StatusOK, response.StatusCode)

	var swaggerDocument struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	suite.Require().NoError(json.NewDecoder(response.Body).Decode(&swaggerDocument))
	suite.Equalf(
		readProjectVersion(suite.T()),
		swaggerDocument.Info.Version,
		"Update the project version in .project-keeper.yml and the Swagger @version annotation in main.go.",
	)
}

func readProjectVersion(t *testing.T) string {
	t.Helper()
	projectKeeperConfiguration, err := os.ReadFile("../../.project-keeper.yml")
	if err != nil {
		t.Fatalf("failed to read the project configuration: %v", err)
	}
	var configuration struct {
		Version string `yaml:"version"`
	}
	if err := yaml.Unmarshal(projectKeeperConfiguration, &configuration); err != nil {
		t.Fatalf("failed to parse the project configuration: %v", err)
	}
	return configuration.Version
}

<<<<<<< HEAD
// [itest->dsn~execute-statement-endpoint~1]
// [itest->dsn~execute-statement-request-body~1]
// [itest->dsn~execute-statement-response-body~1]
func (suite *DockerImageTestSuite) TestExecuteStatementWithMalformedJSONDocker() {
	apiContainer := runRestAPIContainer(suite.restAPIProperties(), suite.exasolPort)
	ip, err := apiContainer.Host(suite.ctx)
	onError(err)
	port, err := apiContainer.MappedPort(suite.ctx, "8080")
	onError(err)

	request, err := http.NewRequest(http.MethodPost, "http://"+ip+":"+port.Port()+"/api/v1/statement", bytes.NewBufferString("{"))
	onError(err)
	request.Header.Set("Authorization", "3J90XAv9loMIXzQdfYmtJrHAbopPsc")
	response, err := http.DefaultClient.Do(request)
	onError(err)
	body, err := io.ReadAll(response.Body)
	onError(err)

	suite.Equal("400 Bad Request", response.Status)
	suite.Equal("{\"status\":\"error\",\"exception\":\"unexpected EOF\"}", string(body))
}

=======
>>>>>>> e36c59f (Verify API version in swagger doc is up-to-date)
// [itest->dsn~execute-statement-endpoint~1]
// [itest->dsn~execute-statement-request-body~1]
// [itest->dsn~execute-statement-response-body~1]
func (suite *DockerImageTestSuite) TestExecuteStatementWithMalformedJSONDocker() {
	apiContainer := runRestAPIContainer(suite.restAPIProperties(), suite.exasolPort)
	ip, err := apiContainer.Host(suite.ctx)
	onError(err)
	port, err := apiContainer.MappedPort(suite.ctx, "8080")
	onError(err)

	request, err := http.NewRequest(http.MethodPost, "http://"+ip+":"+port.Port()+"/api/v1/statement", bytes.NewBufferString("{"))
	onError(err)
	request.Header.Set("Authorization", "3J90XAv9loMIXzQdfYmtJrHAbopPsc")
	response, err := http.DefaultClient.Do(request)
	onError(err)
	body, err := io.ReadAll(response.Body)
	onError(err)

	suite.Equal("400 Bad Request", response.Status)
	suite.Equal("{\"status\":\"error\",\"exception\":\"unexpected EOF\"}", string(body))
}

func (suite *DockerImageTestSuite) restAPIProperties() map[string]string {
	return map[string]string{
		exasol_rest_api.APITokensKey:                       suite.defaultAuthTokens,
		exasol_rest_api.ExasolUserKey:                      suite.defaultExasolUsername,
		exasol_rest_api.ExasolPasswordKey:                  suite.defaultExasolPassword,
		exasol_rest_api.ExasolHostKey:                      testcontainers.HostInternal,
		exasol_rest_api.ExasolPortKey:                      strconv.Itoa(suite.exasolPort),
		exasol_rest_api.ExasolValidateServerCertificateKey: "false",
	}
}

func runRestAPIContainer(env map[string]string, hostAccessPort int) testcontainers.Container {
	image := "rest-api-test-image:latest"
	request := testcontainers.ContainerRequest{
		Image:           image,
		ExposedPorts:    []string{"8080"},
		HostAccessPorts: []int{hostAccessPort},
		WaitingFor:      wait.ForLog("Listening and serving HTTP"),
		Env:             env,
	}
	apiContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("Starting of docker image %q failed with error %q. Run 'docker buildx build --tag %s .' before starting the tests", image, err.Error(), image))
	}
	return apiContainer
}
