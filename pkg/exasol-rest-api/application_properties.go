package exasol_rest_api

import (
	"os"

	exaerror "github.com/exasol/error-reporting-go"
)

const APITokensKey string = "API_TOKENS"
const ApplicationServerKey string = "SERVER_ADDRESS"
const ExasolUserKey string = "EXASOL_USER"
const ExasolPasswordKey string = "EXASOL_PASSWORD"
const ExasolHostKey string = "EXASOL_HOST"
const ExasolPortKey string = "EXASOL_PORT"
const EncryptionKey string = "EXASOL_ENCRYPTION"
const UseTLSKey string = "EXASOL_TLS"

// ApplicationProperties for Exasol REST API service.
// [impl->dsn~service-account~1]
// [impl->dsn~service-credentials~1]
type ApplicationProperties struct {
	APITokens                    []string `yaml:"API_TOKENS"`
	ApplicationServer            string   `yaml:"SERVER_ADDRESS"`
	ExasolUser                   string   `yaml:"EXASOL_USER"`
	ExasolPassword               string   `yaml:"EXASOL_PASSWORD"`
	ExasolHost                   string   `yaml:"EXASOL_HOST"`
	ExasolPort                   int      `yaml:"EXASOL_PORT"`
	Encryption                   int      `yaml:"EXASOL_ENCRYPTION"`
	ExasolCertificateFingerprint string   `yaml:"EXASOL_CERTIFICATE_FINGERPRINT"`
	UseTLS                       int      `yaml:"EXASOL_TLS"`
	APIUseTLS                    bool     `yaml:"API_TLS"`
	APITLSPrivateKeyPath         string   `yaml:"API_TLS_PKPATH"`
	APITLSCertificatePath        string   `yaml:"API_TLS_CERTPATH"`
}

// GetApplicationProperties creates an application properties.
func GetApplicationProperties(appPropertiesPathCli string) *ApplicationProperties {
	properties := readApplicationProperties(appPropertiesPathCli)
	err := properties.validate()
	if err != nil {
		panic(exaerror.New("E-ERA-7").Message("application properties validation failed. {{error|uq}}").
			Parameter("error", err.Error()).Error())
	}
	return &properties
}

func readApplicationProperties(appPropertiesPathCli string) ApplicationProperties {
	properties := readApplicationPropertiesFromFile(appPropertiesPathCli)
	properties.setValuesFromEnvironmentVariables()
	properties.fillMissingWithDefaultValues()
	return properties
}

func readApplicationPropertiesFromFile(appPropertiesPathCli string) ApplicationProperties {
	var propertiesFilePath string
	if appPropertiesPathCli != "" {
		propertiesFilePath = appPropertiesPathCli
	} else {
		propertiesFilePath = os.Getenv("APPLICATION_PROPERTIES_PATH")
	}
	properties, err := getPropertiesFromFile(propertiesFilePath)
	if err != nil {
		errorLogger.Print(exaerror.New("E-ERA-6").
			Message("cannot read properties from specified file: {{file path}}. {{error|uq}}").
			Parameter("file path", propertiesFilePath).Parameter("error", err.Error()).String())
		return ApplicationProperties{}
	} else {
		return properties
	}
}

func (applicationProperties *ApplicationProperties) fillMissingWithDefaultValues() {
	defaultProperties := getDefaultProperties()
	if applicationProperties.ApplicationServer == "" {
		applicationProperties.ApplicationServer = defaultProperties.ApplicationServer
	}
	if applicationProperties.ExasolHost == "" {
		applicationProperties.ExasolHost = defaultProperties.ExasolHost
	}
	if applicationProperties.ExasolPort == 0 {
		applicationProperties.ExasolPort = defaultProperties.ExasolPort
	}
	if applicationProperties.Encryption != -1 && applicationProperties.Encryption != 1 {
		applicationProperties.Encryption = defaultProperties.Encryption
	}
	if applicationProperties.UseTLS != -1 && applicationProperties.UseTLS != 1 {
		applicationProperties.UseTLS = defaultProperties.UseTLS
	}
}

func (applicationProperties *ApplicationProperties) validate() error {
	if applicationProperties.ExasolUser == "" && applicationProperties.ExasolPassword == "" {
		return exaerror.New("E-ERA-8").
			Message("exasol username and password are missing in properties.").
			Mitigation("please specify an Exasol username and password via properties.")
	} else if applicationProperties.ExasolUser == "" {
		return exaerror.New("E-ERA-9").
			Message("exasol username is missing in properties.").
			Mitigation("please specify an Exasol username via properties.")
	} else if applicationProperties.ExasolPassword == "" {
		return exaerror.New("E-ERA-10").
			Message("exasol password is missing in properties.").
			Mitigation("please specify an Exasol password via properties.")
	} else {
		return nil
	}
}

func getDefaultProperties() *ApplicationProperties {
	return &ApplicationProperties{
		ApplicationServer: "0.0.0.0:8080",
		ExasolHost:        "localhost",
		ExasolPort:        8563,
		Encryption:        1,
		UseTLS:            1,
	}
}
