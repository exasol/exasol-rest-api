package exasol_rest_api

import "testing"

// [utest->dsn~get-rows-request-parameters~1]
func TestGetValueByTypeParsesFloatValue(t *testing.T) {
	value, err := getValueByType("float", "15.5")
	if err != nil {
		t.Fatalf("getValueByType() returned an error: %v", err)
	}
	if value != 15.5 {
		t.Errorf("getValueByType() = %v, want 15.5", value)
	}
}
