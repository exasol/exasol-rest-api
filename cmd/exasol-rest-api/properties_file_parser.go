package exasol_rest_api

import (
	"gopkg.in/yaml.v2"
	"os"
)

func getPropertiesFromFile(filepath string, properties *interface{}) error {
	configFile, err := openFile(filepath)
	if err != nil {
		return err
	}
	decodePropertiesFile(configFile, properties)
	closeFile(configFile)
	return nil
}

func openFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	} else if file == nil {
		return nil, os.ErrNotExist
	} else {
		return file, err
	}
}

func decodePropertiesFile(configFile *os.File, properties *interface{}) {
	decoder := yaml.NewDecoder(configFile)
	err := decoder.Decode(&properties)
	if err != nil {
		errorLogger.Printf("cannot decode a property file: %s. %s", configFile.Name(), err)
	}
}

func closeFile(configFile *os.File) {
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			errorLogger.Printf("error closing a file: %s. %s", configFile.Name(), err)
		}
	}(configFile)
}
