package exasol_rest_api

type Command struct {
	Command string `json:"command"`
}

type LoginCommand struct {
	Command
	ProtocolVersion int `json:"protocolVersion"`
}

type Attributes struct {
	ResultSetMaxRows int `json:"resultSetMaxRows,omitempty"`
}

type AuthCommand struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	UseCompression   bool   `json:"useCompression"`
	ClientName       string `json:"clientName,omitempty"`
	ClientOs         string `json:"clientOs,omitempty"`
	ClientOsUsername string `json:"clientOsUsername,omitempty"`
	ClientRuntime    string `json:"clientRuntime,omitempty"`
}

type SQLCommand struct {
	Command
	SQLText    string     `json:"sqlText"`
	Attributes Attributes `json:"attributes,omitempty"`
}
