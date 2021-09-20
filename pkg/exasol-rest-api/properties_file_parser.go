package exasol_rest_api

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

func getPropertiesFromFile(filepath string, properties *ApplicationProperties) error {
	propertiesFile, err := openFile(filepath)
	if err != nil {
		return err
	}
	err = decodePropertiesFile(propertiesFile, properties)
	if err != nil {
		return err
	}
	closeFile(propertiesFile)
	return nil
}

func openFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	} else if file == nil {
		return nil, errors.New("properties file doesn't exist")
	} else {
		return file, nil
	}
}

func decodePropertiesFile(propertiesFile *os.File, properties *ApplicationProperties) error {
	decoder := yaml.NewDecoder(propertiesFile)
	return decoder.Decode(&properties)
}

func closeFile(configFile *os.File) {
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			errorLogger.Printf("error closing a file: %s. %s", configFile.Name(), err)
		}
	}(configFile)
}
