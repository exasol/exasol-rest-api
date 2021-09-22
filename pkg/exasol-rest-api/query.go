package exasol_rest_api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//Application represents the REST API service.
type Application struct {
	Properties *ApplicationProperties
}

//Query allows users to send a query to the Exasol database.
func (application *Application) Query(context *gin.Context) {
	response, err := application.queryExasol(context.Param("query"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"ErrorCode": "EXA-REST-API-1",
			"Message":   err.Error(),
		})
	} else {
		context.Data(http.StatusOK, "application/json", response)
	}
}

func (application *Application) queryExasol(query string) ([]byte, error) {
	connection, err := application.openConnection()
	if err != nil {
		return nil, err
	}
	response, err := connection.executeQuery(query)
	if err != nil {
		return nil, err
	}
	err = connection.close()
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (application *Application) openConnection() (*websocketConnection, error) {
	connection := &websocketConnection{
		connProperties: application.Properties,
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
