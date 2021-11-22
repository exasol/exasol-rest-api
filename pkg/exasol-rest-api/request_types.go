package exasol_rest_api

type command struct {
	Command string `json:"command"`
}

type loginCommand struct {
	command
	ProtocolVersion int `json:"protocolVersion"`
}

type attributes struct {
	Autocommit       bool `json:"autocommit,omitempty"`
	ResultSetMaxRows int  `json:"resultSetMaxRows,omitempty"`
}

type authCommand struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	UseCompression   bool   `json:"useCompression"`
	ClientName       string `json:"clientName,omitempty"`
	ClientOs         string `json:"clientOs,omitempty"`
	ClientOsUsername string `json:"clientOsUsername,omitempty"`
	ClientRuntime    string `json:"clientRuntime,omitempty"`
}

type sqlCommand struct {
	command
	SQLText    string     `json:"sqlText"`
	Attributes attributes `json:"attributes,omitempty"`
}
