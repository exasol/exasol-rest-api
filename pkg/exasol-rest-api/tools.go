package exasol_rest_api

// CreateStringsSet converts a slice of string into a set-like map.
func CreateStringsSet(tokens []string) map[string]bool {
	set := make(map[string]bool)
	for _, token := range tokens {
		set[token] = true
	}
	return set
}
