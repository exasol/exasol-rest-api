package exasol_rest_api

import (
	"fmt"
	"os"

	exaerror "github.com/exasol/error-reporting-go"
	"gopkg.in/yaml.v3"
)

func getPropertiesFromFile(filepath string) (ApplicationProperties, error) {
	propertiesFile, err := openFile(filepath)
	if err != nil {
		return ApplicationProperties{}, err
	}

	defer closeFile(propertiesFile)
	return decodePropertiesFile(propertiesFile)
}

func openFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, exaerror.New("E-ERA-11").Message("cannot open a file. {{error|uq}}").
			Parameter("error", err.Error())
	} else if file == nil {
		return nil, fmt.Errorf(exaerror.New("E-ERA-12").
			Message("properties file doesn't exist.").String())
	} else {
		return file, nil
	}
}

func decodePropertiesFile(propertiesFile *os.File) (ApplicationProperties, error) {
	decoder := yaml.NewDecoder(propertiesFile)
	properties := ApplicationProperties{}

	err := decoder.Decode(&properties)
	if err != nil {
		return ApplicationProperties{}, exaerror.New("E-ERA-13").Message("cannot decode properties file. {{error|uq}}.").
			Parameter("error", err.Error()).
			Mitigation("Please make sure that file is not empty and contains correct properties.")
	} else {
		return properties, nil
	}
}

func closeFile(configFile *os.File) {
	if err := configFile.Close(); err != nil {
		errorLogger.Printf("error closing a file: %s. %s", configFile.Name(), err)
	}
}
