package main

import (
	"flag"
	exasol_rest_api "main/pkg/exasol-rest-api"
)

// @title Exasol REST API
// @version 0.2.15
// @description This service is a proxy that wrapping up Exasol WebSockets library.

// @contact.name Exasol REST API GitHub Issues
// @contact.url https://github.com/exasol/exasol-rest-api/issues

// @license.name MIT License
// @license.url https://github.com/exasol/exasol-rest-api/blob/main/LICENSE

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /api/v1

func main() {
	appPropertiesPathCli := extractAppPropertiesPath()
	exasol_rest_api.Run(*appPropertiesPathCli)
}

func extractAppPropertiesPath() *string {
	appPropertiesPathCli := flag.String("application-properties-path", "", "Path to the application properties file.")
	flag.Parse()
	return appPropertiesPathCli
}
