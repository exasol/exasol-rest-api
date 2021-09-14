package exasol_rest_api

import (
	"github.com/mitchellh/mapstructure"
)

func Query(query string) (string, error) {
	connection, _ := openConnection()
	response, _ := connection.executeQuery(query)
	err := connection.close()
	if err != nil {
		return "", err
	}
	return response, nil
}

func openConnection() (*websocketConnection, error) {
	connProperties := readConnectionProperties()
	connection := &websocketConnection{
		connProperties: &connProperties,
	}
	err := connection.connect()
	if err != nil {
		return nil, err
	}
	err = connection.login()
	if err != nil {
		return nil, err
	}
	return connection, err
}

func readConnectionProperties() connectionProperties {
	var properties connectionProperties
	var propertiesAsInterface interface{} = properties
	GetPropertiesFromFile("connection-properties.yml", &propertiesAsInterface)
	err := mapstructure.Decode(propertiesAsInterface, &properties)
	if err != nil {
		ErrorLogger.Printf("error reading websocketConnection properties: %s", err)
	}
	return properties
}
