package exasol_rest_api

import "encoding/json"

type BaseResponse struct {
	Status       string          `json:"status"`
	ResponseData json.RawMessage `json:"responseData"`
	Exception    *Exception      `json:"exception"`
}

type Exception struct {
	Text    string `json:"text"`
	SQLCode string `json:"sqlCode"`
}

type PublicKeyResponse struct {
	PublicKeyPem      string `json:"publicKeyPem"`
	PublicKeyModulus  string `json:"publicKeyModulus"`
	PublicKeyExponent string `json:"publicKeyExponent"`
}
