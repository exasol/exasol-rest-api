package exasol_rest_api_test

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"os"
	"testing"
)

type ApplicationPropertiesSuite struct {
	suite.Suite
}

func TestApplicationPropertiesSuite(t *testing.T) {
	suite.Run(t, new(ApplicationPropertiesSuite))
}

func (suite *ApplicationPropertiesSuite) TestReadingProperties() {
	expected := &exasol_rest_api.ApplicationProperties{
		ApplicationServer:         "test:8888",
		ExasolUser:                "myUser",
		ExasolPassword:            "pass",
		ExasolHost:                "127.0.0.1",
		ExasolPort:                1234,
		Encryption:                true,
		UseTLS:                    true,
		ExasolWebsocketAPIVersion: 3,
	}
	applicationPropertiesPathKey := suite.setPathToPropertiesFileEnv(expected)
	actual := exasol_rest_api.GetApplicationProperties(applicationPropertiesPathKey)
	suite.Equal(expected, actual)
}

func (suite *ApplicationPropertiesSuite) TestDefaultProperties() {
	minimalRequiredProperties := &exasol_rest_api.ApplicationProperties{
		ExasolUser:     "myUser",
		ExasolPassword: "pass",
	}
	applicationPropertiesPathKey := suite.setPathToPropertiesFileEnv(minimalRequiredProperties)
	actual := exasol_rest_api.GetApplicationProperties(applicationPropertiesPathKey)
	expected := &exasol_rest_api.ApplicationProperties{
		ApplicationServer:         "localhost:8080",
		ExasolUser:                "myUser",
		ExasolPassword:            "pass",
		ExasolHost:                "localhost",
		ExasolPort:                8563,
		Encryption:                false,
		UseTLS:                    false,
		ExasolWebsocketAPIVersion: 2,
	}
	suite.Equal(expected, actual)
}

func (suite *ApplicationPropertiesSuite) TestReadingPropertiesWithoutPath() {
	suite.PanicsWithValue("E-ERA-4: runtime error: missing environment variable: 'DUMMY_KEY'",
		func() { exasol_rest_api.GetApplicationProperties("DUMMY_KEY") })
}

func (suite *ApplicationPropertiesSuite) TestReadingPropertiesWithMissingFile() {
	applicationPropertiesPathKey := "APPLICATION_PROPERTIES_PATH"
	err := os.Setenv(applicationPropertiesPathKey, "file/does/not/exist")
	onError(err)
	suite.PanicsWithValue("E-ERA-5: runtime error: application properties are missing or incorrect. "+
		"E-ERA-6: cannot read properties from specified file: 'file/does/not/exist'. "+
		"E-ERA-11: cannot open a file. open file/does/not/exist: no such file or directory",
		func() { exasol_rest_api.GetApplicationProperties("APPLICATION_PROPERTIES_PATH") })
}

func (suite *ApplicationPropertiesSuite) TestDefaultPropertiesWithMissingUsername() {
	properties := &exasol_rest_api.ApplicationProperties{
		ExasolPassword: "pass",
	}
	applicationPropertiesPathKey := suite.setPathToPropertiesFileEnv(properties)
	suite.PanicsWithValue("E-ERA-5: runtime error: application properties are missing or incorrect. "+
		"E-ERA-7: properties file validation failed. "+
		"E-ERA-9: exasol username is missing in properties. please specify an Exasol username via properties.",
		func() { exasol_rest_api.GetApplicationProperties(applicationPropertiesPathKey) })
}

func (suite *ApplicationPropertiesSuite) TestDefaultPropertiesWithMissingPassword() {
	properties := &exasol_rest_api.ApplicationProperties{
		ExasolUser: "myUSer",
	}
	applicationPropertiesPathKey := suite.setPathToPropertiesFileEnv(properties)
	suite.PanicsWithValue("E-ERA-5: runtime error: application properties are missing or incorrect. "+
		"E-ERA-7: properties file validation failed. "+
		"E-ERA-10: exasol password is missing in properties. please specify an Exasol password via properties.",
		func() { exasol_rest_api.GetApplicationProperties(applicationPropertiesPathKey) })
}

func (suite *ApplicationPropertiesSuite) TestDefaultPropertiesWithMissingUsernameAndPassword() {
	properties := &exasol_rest_api.ApplicationProperties{
		UseTLS: true,
	}
	applicationPropertiesPathKey := suite.setPathToPropertiesFileEnv(properties)
	suite.PanicsWithValue("E-ERA-5: runtime error: application properties are missing or incorrect. "+
		"E-ERA-7: properties file validation failed. "+
		"E-ERA-8: exasol username and password are missing in properties. "+
		"please specify an Exasol username and password via properties.",
		func() { exasol_rest_api.GetApplicationProperties(applicationPropertiesPathKey) })
}

func (suite *ApplicationPropertiesSuite) setPathToPropertiesFileEnv(
	properties *exasol_rest_api.ApplicationProperties) string {
	file, err := ioutil.TempFile("", "application_properties_*.yml")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	onError(err)
	data, err := yaml.Marshal(&properties)
	onError(err)
	_, err = file.Write(data)
	onError(err)
	applicationPropertiesPathKey := "APPLICATION_PROPERTIES_PATH"
	err = os.Setenv(applicationPropertiesPathKey, file.Name())
	onError(err)
	return applicationPropertiesPathKey
}
