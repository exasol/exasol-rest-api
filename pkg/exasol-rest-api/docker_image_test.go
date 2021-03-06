package exasol_rest_api_test

import (
	"context"
	"fmt"
	"io/ioutil"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DockerImageTestSuite struct {
	suite.Suite
	ctx                   context.Context
	exasolContainer       testcontainers.Container
	defaultExasolUsername string
	defaultExasolPassword string
	defaultAuthTokens     string
	exasolPort            int
	exasolHost            string
	exasolContainerIP     string
}

func TestDockerImageSuite(t *testing.T) {
	suite.Run(t, new(DockerImageTestSuite))
}

func (suite *DockerImageTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.defaultExasolUsername = "api_service_account"
	suite.defaultExasolPassword = "secret_password"
	suite.defaultAuthTokens = "3J90XAv9loMIXzQdfYmtJrHAbopPsc,OR6rq6KjWmhvGU770A9OTjpfH86nlk"
	suite.exasolContainer = runExasolContainer(suite.ctx)
	suite.exasolHost = getExasolHost(suite.exasolContainer, suite.ctx)
	suite.exasolPort = getExasolPort(suite.exasolContainer, suite.ctx)
	suite.exasolContainerIP = getExasolContainerIP(suite.exasolContainer, suite.ctx)
	createDefaultServiceUserWithAccess(suite.defaultExasolUsername, suite.defaultExasolPassword, suite.exasolHost,
		suite.exasolPort)
}

func (suite *DockerImageTestSuite) TestQuery() {
	properties := map[string]string{
		exasol_rest_api.APITokensKey:      suite.defaultAuthTokens,
		exasol_rest_api.ExasolUserKey:     suite.defaultExasolUsername,
		exasol_rest_api.ExasolPasswordKey: suite.defaultExasolPassword,
		exasol_rest_api.ExasolHostKey:     suite.exasolContainerIP,
		exasol_rest_api.ExasolPortKey:     "8563",
		exasol_rest_api.EncryptionKey:     "-1",
	}
	apiContainer := runRestAPIContainer(properties)
	ip, err := apiContainer.Host(context.Background())
	onError(err)

	port, err := apiContainer.MappedPort(suite.ctx, "8080")
	onError(err)

	baseUrl := "http://" + ip + ":" + strconv.Itoa(port.Int())
	req, err := http.NewRequest(http.MethodGet, baseUrl+"/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	req.Header.Set("Authorization", "3J90XAv9loMIXzQdfYmtJrHAbopPsc")
	onError(err)

	client := http.Client{}
	response, err := client.Do(req)
	onError(err)

	body, err := ioutil.ReadAll(response.Body)
	onError(err)

	suite.Equal("200 OK", response.Status)
	suite.Equal("{\"status\":\"ok\",\"rows\":[{\"X\":15,\"Y\":\"test\"},{\"X\":10,\"Y\":\"test_2\"}],\"meta\":{\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100,\"characterSet\":\"UTF8\"}}]}}",
		string(body))
}

func runRestAPIContainer(env map[string]string) testcontainers.Container {
	image := "rest-api-test-image:latest"
	request := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"8080"},
		WaitingFor:   wait.ForLog("Listening and serving HTTP"),
		Env:          env,
	}
	apiContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("Starting of docker image %q failed with error %q. Run 'docker build --tag %s .' before starting the tests", image, err.Error(), image))
	}
	return apiContainer
}
