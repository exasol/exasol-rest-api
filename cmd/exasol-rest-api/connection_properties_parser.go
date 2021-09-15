package exasol_rest_api

type ConnectionProperties struct {
	User             string
	Password         string
	Host             string
	Port             int
	ApiVersion       int
	FetchSize        int
	ResultSetMaxRows int
	Encryption       bool
	UseTLS           bool
}

func createConnectionProperties(userDefinedProperties ConnectionProperties) *ConnectionProperties {
	defaultConfig := getDefaultConfig(userDefinedProperties.Host, userDefinedProperties.Port)
	defaultConfig.User = userDefinedProperties.User
	defaultConfig.Password = userDefinedProperties.Password
	defaultConfig.Port = userDefinedProperties.Port
	defaultConfig.ApiVersion = userDefinedProperties.ApiVersion
	defaultConfig.FetchSize = userDefinedProperties.FetchSize
	defaultConfig.ResultSetMaxRows = userDefinedProperties.ResultSetMaxRows
	defaultConfig.Encryption = userDefinedProperties.Encryption
	defaultConfig.UseTLS = userDefinedProperties.UseTLS
	return defaultConfig
}

func getDefaultConfig(host string, port int) *ConnectionProperties {
	return &ConnectionProperties{
		Host:       host,
		Port:       port,
		ApiVersion: 2,
		Encryption: true,
		UseTLS:     false,
		FetchSize:  128 * 1024,
	}
}
