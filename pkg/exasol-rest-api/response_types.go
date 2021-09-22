package exasol_rest_api

import "encoding/json"

type baseResponse struct {
	Status       string          `json:"status"`
	ResponseData json.RawMessage `json:"responseData"`
	Exception    *exception      `json:"exception"`
}

type exception struct {
	Text    string `json:"text"`
	SQLCode string `json:"sqlCode"`
}

type publicKeyResponse struct {
	PublicKeyPem      string `json:"publicKeyPem"`
	PublicKeyModulus  string `json:"publicKeyModulus"`
	PublicKeyExponent string `json:"publicKeyExponent"`
}
