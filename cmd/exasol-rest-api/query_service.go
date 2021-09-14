package exasol_rest_api

import (
	"context"
	"github.com/mitchellh/mapstructure"
)

func Query(query string) (string, error) {
	connection, _ := openConnection()
	response, _ := connection.simpleExec(query)
	err := connection.Close()
	if err != nil {
		return "", err
	}
	return response, nil
}

func openConnection() (*connection, error) {
	connProperties := readConnectionProperties()
	exasolConnector := &connector{
		connProperties: createConnectionProperties(connProperties),
	}
	return exasolConnector.Connect(context.Background())
}

func readConnectionProperties() connectionProperties {
	var properties connectionProperties
	var propertiesAsInterface interface{} = properties
	GetPropertiesFromFile("connection-properties.yml", &propertiesAsInterface)
	err := mapstructure.Decode(propertiesAsInterface, &properties)
	if err != nil {
		ErrorLogger.Printf("error reading connection properties: %s", err)
	}
	return properties
}
