package main

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

type AuthResponse struct {
	SessionID             int    `json:"sessionId"`
	ProtocolVersion       int    `json:"protocolVersion"`
	ReleaseVersion        string `json:"releaseVersion"`
	DatabaseName          string `json:"databaseName"`
	ProductName           string `json:"productName"`
	MaxDataMessageSize    int    `json:"maxDataMessageSize"`
	MaxIdentifierLength   int    `json:"maxIdentifierLength"`
	MaxVarcharLength      int    `json:"maxVarcharLength"`
	IdentifierQuoteString string `json:"identifierQuoteString"`
	TimeZone              string `json:"timeZone"`
	TimeZoneBehavior      string `json:"timeZoneBehavior"`
}

type PublicKeyResponse struct {
	PublicKeyPem      string `json:"publicKeyPem"`
	PublicKeyModulus  string `json:"publicKeyModulus"`
	PublicKeyExponent string `json:"publicKeyExponent"`
}
