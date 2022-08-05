package exasol_rest_api_test

import (
	"io/ioutil"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

const applicationPropertiesPathKey = "APPLICATION_PROPERTIES_PATH"

type ApplicationPropertiesSuite struct {
	suite.Suite
}

func TestApplicationPropertiesSuite(t *testing.T) {
	suite.Run(t, new(ApplicationPropertiesSuite))
}

func (suite *ApplicationPropertiesSuite) TestReadingProperties() {
	expected := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{"abc"},
		ApplicationServer:         "test:8888",
		ExasolUser:                "myUser",
		ExasolPassword:            "pass",
		ExasolHost:                "127.0.0.1",
		ExasolPort:                1234,
		Encryption:                1,
		UseTLS:                    1,
		ExasolWebsocketAPIVersion: 3,
	}
	suite.setPathToPropertiesFileEnv(expected)
	actual := exasol_rest_api.GetApplicationProperties("")
	suite.Equal(expected, actual)
}
func (suite *ApplicationPropertiesSuite) TestViaCLIArgumentParam() {
	expected := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{"abc"},
		ApplicationServer:         "test:8888",
		ExasolUser:                "myUser",
		ExasolPassword:            "pass",
		ExasolHost:                "127.0.0.1",
		ExasolPort:                1234,
		Encryption:                1,
		UseTLS:                    1,
		ExasolWebsocketAPIVersion: 3,
	}
	file, err := suite.createTempAppPropertiesFile("", "application_properties_cli*.yml", expected)
	onError(err)
	actual := exasol_rest_api.GetApplicationProperties(file.Name())
	suite.Equal(expected, actual)
}
func (suite *ApplicationPropertiesSuite) TestDefaultProperties() {
	minimalRequiredProperties := &exasol_rest_api.ApplicationProperties{
		ExasolUser:     "myUser",
		ExasolPassword: "pass",
	}
	suite.setPathToPropertiesFileEnv(minimalRequiredProperties)
	actual := exasol_rest_api.GetApplicationProperties("")
	expected := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{},
		ApplicationServer:         "0.0.0.0:8080",
		ExasolUser:                "myUser",
		ExasolPassword:            "pass",
		ExasolHost:                "localhost",
		ExasolPort:                8563,
		Encryption:                1,
		UseTLS:                    1,
		ExasolWebsocketAPIVersion: 2,
	}
	suite.Equal(expected, actual)
}

func (suite *ApplicationPropertiesSuite) TestReadingPropertiesWithMissingPropertiesFileAndWithoutEnv() {
	err := os.Setenv(applicationPropertiesPathKey, "file/does/not/exist")
	onError(err)
	suite.PanicsWithValue("E-ERA-7: application properties validation failed. "+
		"E-ERA-8: exasol username and password are missing in properties. "+
		"please specify an Exasol username and password via properties.",
		func() { exasol_rest_api.GetApplicationProperties("") })
}

func (suite *ApplicationPropertiesSuite) TestReadingPropertiesWithEmptyPropertiesFileAndWithoutEnv() {
	file, _ := ioutil.TempFile("", "application_properties_*.yml")
	defer func(file *os.File) {
		onError(file.Close())
	}(file)

	err := os.Setenv(applicationPropertiesPathKey, file.Name())
	onError(err)
	suite.PanicsWithValue("E-ERA-7: application properties validation failed. "+
		"E-ERA-8: exasol username and password are missing in properties. "+
		"please specify an Exasol username and password via properties.",
		func() { exasol_rest_api.GetApplicationProperties("") })
}

func (suite *ApplicationPropertiesSuite) TestDefaultPropertiesWithMissingUsername() {
	properties := &exasol_rest_api.ApplicationProperties{
		ExasolPassword: "pass",
	}
	suite.setPathToPropertiesFileEnv(properties)
	suite.PanicsWithValue("E-ERA-7: application properties validation failed. "+
		"E-ERA-9: exasol username is missing in properties. please specify an Exasol username via properties.",
		func() { exasol_rest_api.GetApplicationProperties("") })
}

func (suite *ApplicationPropertiesSuite) TestDefaultPropertiesWithMissingPassword() {
	properties := &exasol_rest_api.ApplicationProperties{
		ExasolUser: "myUSer",
	}
	suite.setPathToPropertiesFileEnv(properties)
	suite.PanicsWithValue("E-ERA-7: application properties validation failed. "+
		"E-ERA-10: exasol password is missing in properties. please specify an Exasol password via properties.",
		func() { exasol_rest_api.GetApplicationProperties("") })
}

func (suite *ApplicationPropertiesSuite) TestDefaultPropertiesWithMissingUsernameAndPassword() {
	properties := &exasol_rest_api.ApplicationProperties{
		UseTLS: 1,
	}
	suite.setPathToPropertiesFileEnv(properties)
	suite.PanicsWithValue("E-ERA-7: application properties validation failed. "+
		"E-ERA-8: exasol username and password are missing in properties. "+
		"please specify an Exasol username and password via properties.",
		func() { exasol_rest_api.GetApplicationProperties("") })
}

func (suite *ApplicationPropertiesSuite) TestReadingPropertiesWithEnv() {
	expected := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{"abc", "bca"},
		ApplicationServer:         "test:8888",
		ExasolUser:                "myUser",
		ExasolPassword:            "pass",
		ExasolHost:                "127.0.0.1",
		ExasolPort:                1234,
		Encryption:                1,
		UseTLS:                    -1,
		ExasolWebsocketAPIVersion: 3,
	}
	err := os.Setenv(exasol_rest_api.APITokensKey, "abc,bca")
	onError(err)
	err = os.Setenv(exasol_rest_api.ApplicationServerKey, "test:8888")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolUserKey, "myUser")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolPasswordKey, "pass")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolHostKey, "127.0.0.1")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolPortKey, "1234")
	onError(err)
	err = os.Setenv(exasol_rest_api.EncryptionKey, "1")
	onError(err)
	err = os.Setenv(exasol_rest_api.UseTLSKey, "-1")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolWebsocketAPIVersionKey, "3")
	onError(err)
	actual := exasol_rest_api.GetApplicationProperties("")
	suite.Equal(expected, actual)
}

func (suite *ApplicationPropertiesSuite) TestOverridingPropertiesFromFileWithEnv() {
	propertiesFromFile := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{"abc"},
		ApplicationServer:         "1.1.1.1:8888",
		ExasolUser:                "user",
		ExasolPassword:            "pass111",
		ExasolHost:                "localhost1",
		ExasolPort:                4321,
		Encryption:                1,
		UseTLS:                    1,
		ExasolWebsocketAPIVersion: 2,
	}
	suite.setPathToPropertiesFileEnv(propertiesFromFile)
	expected := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{"abc", "bca"},
		ApplicationServer:         "test:8888",
		ExasolUser:                "myUser",
		ExasolPassword:            "pass",
		ExasolHost:                "127.0.0.1",
		ExasolPort:                1234,
		Encryption:                -1,
		UseTLS:                    -1,
		ExasolWebsocketAPIVersion: 3,
	}
	err := os.Setenv(exasol_rest_api.APITokensKey, "abc,bca")
	onError(err)
	err = os.Setenv(exasol_rest_api.ApplicationServerKey, "test:8888")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolUserKey, "myUser")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolPasswordKey, "pass")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolHostKey, "127.0.0.1")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolPortKey, "1234")
	onError(err)
	err = os.Setenv(exasol_rest_api.EncryptionKey, "-1")
	onError(err)
	err = os.Setenv(exasol_rest_api.UseTLSKey, "-1")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolWebsocketAPIVersionKey, "3")
	onError(err)
	actual := exasol_rest_api.GetApplicationProperties("")
	suite.Equal(expected, actual)
}

func (suite *ApplicationPropertiesSuite) TestMixingPropertiesFromFileAndEnv() {
	propertiesFromFile := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{"abc"},
		ApplicationServer:         "1.1.1.1:8888",
		ExasolUser:                "user",
		ExasolPassword:            "pass111",
		ExasolHost:                "localhost1",
		ExasolPort:                4321,
		Encryption:                -1,
		UseTLS:                    -1,
		ExasolWebsocketAPIVersion: 2,
	}
	suite.setPathToPropertiesFileEnv(propertiesFromFile)
	expected := &exasol_rest_api.ApplicationProperties{
		APITokens:                 []string{"abc", "bca"},
		ApplicationServer:         "test:8888",
		ExasolUser:                "user",
		ExasolPassword:            "pass",
		ExasolHost:                "127.0.0.1",
		ExasolPort:                4321,
		Encryption:                -1,
		UseTLS:                    1,
		ExasolWebsocketAPIVersion: 2,
	}
	err := os.Setenv(exasol_rest_api.APITokensKey, "abc,bca")
	onError(err)
	err = os.Setenv(exasol_rest_api.ApplicationServerKey, "test:8888")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolPasswordKey, "pass")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolHostKey, "127.0.0.1")
	onError(err)
	err = os.Setenv(exasol_rest_api.ExasolPortKey, "wrong")
	onError(err)
	err = os.Setenv(exasol_rest_api.EncryptionKey, "bad")
	onError(err)
	err = os.Setenv(exasol_rest_api.UseTLSKey, "1")
	onError(err)
	actual := exasol_rest_api.GetApplicationProperties("")
	suite.Equal(expected, actual)
}

func (suite *ApplicationPropertiesSuite) setPathToPropertiesFileEnv(
	properties *exasol_rest_api.ApplicationProperties) string {
	var filePath string = ""
	var fileNamePattern string = "application_properties_*.yml"
	file, err := suite.createTempAppPropertiesFile(filePath, fileNamePattern, properties)
	onError(err)
	err = os.Setenv(applicationPropertiesPathKey, file.Name())
	onError(err)
	return applicationPropertiesPathKey
}

func (suite *ApplicationPropertiesSuite) createTempAppPropertiesFile(filePath string, fileNamePattern string, properties *exasol_rest_api.ApplicationProperties) (*os.File, error) {
	file, _ := ioutil.TempFile(filePath, fileNamePattern)
	defer func(file *os.File) {
		onError(file.Close())
	}(file)

	data, err := yaml.Marshal(&properties)
	onError(err)
	_, err = file.Write(data)
	onError(err)
	return file, err
}

func (suite *ApplicationPropertiesSuite) SetupTest() {
	err := os.Unsetenv(exasol_rest_api.APITokensKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.ApplicationServerKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.ExasolUserKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.ExasolPasswordKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.ExasolHostKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.ExasolPortKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.EncryptionKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.UseTLSKey)
	onError(err)
	err = os.Unsetenv(exasol_rest_api.ExasolWebsocketAPIVersionKey)
	onError(err)
	err = os.Unsetenv(applicationPropertiesPathKey)
	onError(err)
}
