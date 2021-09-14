package main

import "context"

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
	config := getPropertiesFromFile()
	exasolConnector := &connector{
		config: config,
	}
	return exasolConnector.Connect(context.Background())
}
