package exasol_rest_api_test

import (
	"context"
	"fmt"
	"io/ioutil"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"net/http"
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
	suite.exasolPort = 8563
	createDefaultServiceUserWithAccess(suite.defaultExasolUsername, suite.defaultExasolPassword, suite.exasolHost,
		suite.exasolPort)
}

func (suite *DockerImageTestSuite) TestQuery() {
	properties := map[string]string{
		exasol_rest_api.APITokensKey:      suite.defaultAuthTokens,
		exasol_rest_api.ExasolUserKey:     suite.defaultExasolUsername,
		exasol_rest_api.ExasolPasswordKey: suite.defaultExasolPassword,
		exasol_rest_api.ExasolHostKey:     suite.exasolHost,
	}
	apiContainer := runRestAPIContainer(properties)
	ip, err := apiContainer.ContainerIP(context.Background())
	onError(err)

	req, err := http.NewRequest(http.MethodGet,
		"http://"+ip+":8080/api/v1/query/SELECT * FROM TEST_SCHEMA_1.TEST_TABLE", nil)
	req.Header.Set("Authorization", "3J90XAv9loMIXzQdfYmtJrHAbopPsc")
	onError(err)

	client := http.Client{}
	response, err := client.Do(req)
	onError(err)

	body, err := ioutil.ReadAll(response.Body)
	onError(err)

	suite.Equal("200 OK", response.Status)
	suite.Equal("{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":1,\"numRowsInMessage\":1,\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100,\"characterSet\":\"UTF8\"}}],\"data\":[[15],[\"test\"]]}}],\"numResults\":1}}",
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
		panic(fmt.Errorf("Failed to start docker image %q. Run 'docker build --tag %s .'", image, image))
	}
	return apiContainer
}
