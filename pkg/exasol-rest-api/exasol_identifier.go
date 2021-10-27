package exasol_rest_api

import "strings"

// ToExasolIdentifier makes an Exasol Identifier from a string.
func ToExasolIdentifier(identifier string) string {
	return "\"" + strings.ReplaceAll(identifier, "\"", "\"\"") + "\""
}
