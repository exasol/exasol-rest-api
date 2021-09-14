package exasol_rest_api

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

func Query(context *gin.Context) {
	response, err := QueryExasol(context.Param("query"))
	if err != nil {
		ErrorLogger.Printf("error during querying Exasol: %s", err)
	} else {
		context.IndentedJSON(http.StatusOK, response)
	}
}

func QueryExasol(query string) (string, error) {
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
		connProperties: connProperties,
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

func readConnectionProperties() *connectionProperties {
	var properties connectionProperties
	var propertiesAsInterface interface{} = properties
	err := GetPropertiesFromFile("/home/sia/git/exasol-rest-api/cmd/exasol-rest-api/connection-properties.yml", &propertiesAsInterface)
	if err != nil {
		return nil
	}
	err = mapstructure.Decode(propertiesAsInterface, &properties)
	if err != nil {
		ErrorLogger.Printf("error reading websocketConnection properties: %s", err)
	}
	return createConnectionProperties(properties)
}
