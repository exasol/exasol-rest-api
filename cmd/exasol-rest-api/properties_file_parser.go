package exasol_rest_api

import (
	"gopkg.in/yaml.v2"
	"os"
)

func GetPropertiesFromFile(filepath string, properties *interface{}) {
	configFile := openFile(filepath)
	decodePropertiesFile(configFile, properties)
	closeFile(configFile)
}

func openFile(filepath string) *os.File {
	configFile, err := os.Open(filepath)
	if err != nil {
		ErrorLogger.Printf("cannot open a file: %s. %s", filepath, err)
	}
	return configFile
}

func decodePropertiesFile(configFile *os.File, properties *interface{}) {
	decoder := yaml.NewDecoder(configFile)
	err := decoder.Decode(&properties)
	if err != nil {
		ErrorLogger.Printf("cannot decode a property file: %s. %s", configFile.Name(), err)
	}
}

func closeFile(configFile *os.File) {
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			ErrorLogger.Printf("error closing a file: %s. %s", configFile.Name(), err)
		}
	}(configFile)
}
