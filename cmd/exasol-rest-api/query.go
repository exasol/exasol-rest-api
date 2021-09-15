package exasol_rest_api

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"os"
)

func Query(context *gin.Context) {
	response, err := QueryExasol(context.Param("query"))
	if err != nil {
		errorLogger.Printf("error during querying Exasol: %s", err)
	} else {
		context.IndentedJSON(http.StatusOK, response)
	}
}

func QueryExasol(query string) (string, error) {
	propertiesPath := os.Getenv("CONNECTION_PROPERTIES_PATH")
	connection, _ := openConnection(propertiesPath)
	response, _ := connection.executeQuery(query)
	err := connection.close()
	if err != nil {
		return "", err
	}
	return response, nil
}

func openConnection(propertiesPath string) (*websocketConnection, error) {
	connProperties := readConnectionProperties(propertiesPath)
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

func readConnectionProperties(propertiesPath string) *ConnectionProperties {
	var properties ConnectionProperties
	var propertiesAsInterface interface{} = properties
	err := getPropertiesFromFile(propertiesPath, &propertiesAsInterface)
	if err != nil {
		return nil
	}
	err = mapstructure.Decode(propertiesAsInterface, &properties)
	if err != nil {
		errorLogger.Printf("error reading websocketConnection properties: %s", err)
	}
	return createConnectionProperties(properties)
}
