package exasol_rest_api

import "strings"

// toExasolIdentifier makes an Exasol Identifier from a string.
func toExasolIdentifier(identifier string) string {
	return "\"" + strings.ReplaceAll(identifier, "\"", "\"\"") + "\""
}
