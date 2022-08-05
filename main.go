package main

import (
	"flag"
	exasol_rest_api "main/pkg/exasol-rest-api"
)

// @title Exasol REST API
// @version 0.2.4
// @description This service is a proxy that uses Exasol WebSockets library.

// @contact.name Exasol REST API GitHub Issues
// @contact.url https://github.com/exasol/exasol-rest-api/issues

// @license.name MIT License
// @license.url https://github.com/exasol/exasol-rest-api/blob/main/LICENSE

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /api/v1

func main() {
	app_properties_path := extractAppPropertiesPath()
	exasol_rest_api.Run(*app_properties_path)
}

func extractAppPropertiesPath() *string {
	app_properties_path := flag.String("application-properties-path", "", "Option to provide the application properties path via CLI argument.")
	flag.Parse()
	return app_properties_path
}
