package exasol_rest_api

import (
	"time"
)

type connectionProperties struct {
	User             string
	Password         string
	Host             string
	Port             int
	Params           map[string]string // Connection parameters
	ApiVersion       int
	ClientName       string
	ClientVersion    string
	Schema           string
	FetchSize        int
	Compression      bool
	ResultSetMaxRows int
	Timeout          time.Time
	Encryption       bool
	UseTLS           bool
}

func createConnectionProperties(userDefinedProperties connectionProperties) *connectionProperties {
	defaultConfig := getDefaultConfig(userDefinedProperties.Host, userDefinedProperties.Port)
	defaultConfig.User = userDefinedProperties.User
	defaultConfig.Password = userDefinedProperties.Password
	defaultConfig.Port = userDefinedProperties.Port
	defaultConfig.Params = userDefinedProperties.Params
	defaultConfig.ApiVersion = userDefinedProperties.ApiVersion
	defaultConfig.ClientName = userDefinedProperties.ClientName
	defaultConfig.ClientVersion = userDefinedProperties.ClientVersion
	defaultConfig.Schema = userDefinedProperties.Schema
	defaultConfig.FetchSize = userDefinedProperties.FetchSize
	defaultConfig.Compression = userDefinedProperties.Compression
	defaultConfig.ResultSetMaxRows = userDefinedProperties.ResultSetMaxRows
	defaultConfig.Timeout = userDefinedProperties.Timeout
	defaultConfig.Encryption = userDefinedProperties.Encryption
	defaultConfig.UseTLS = userDefinedProperties.UseTLS
	return defaultConfig
}

func getDefaultConfig(host string, port int) *connectionProperties {
	return &connectionProperties{
		Host:        host,
		Port:        port,
		ApiVersion:  2,
		Encryption:  true,
		Compression: false,
		UseTLS:      false,
		ClientName:  "Go client",
		Params:      map[string]string{},
		FetchSize:   128 * 1024,
	}
}
