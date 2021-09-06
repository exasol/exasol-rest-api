package main

type Command struct {
	Command string `json:"command"`
}

type LoginCommand struct {
	Command
	ProtocolVersion int        `json:"protocolVersion"`
	Attributes      Attributes `json:"attributes,omitempty"`
}

type Attributes struct {
	Autocommit                  bool   `json:"autocommit,omitempty"`
	CompressionEnabled          bool   `json:"compressionEnabled,omitempty"`
	CurrentSchema               string `json:"currentSchema,omitempty"`
	DateFormat                  string `json:"dateFormat,omitempty"`
	DateLanguage                string `json:"dateLanguage,omitempty"`
	DatetimeFormat              string `json:"datetimeFormat,omitempty"`
	DefaultLikeEscapeCharacter  string `json:"defaultLikeEscapeCharacter,omitempty"`
	FeedbackInterval            int    `json:"feedbackInterval,omitempty"`
	NumericCharacters           string `json:"numericCharacters,omitempty"`
	OpenTransaction             bool   `json:"openTransaction,omitempty"`
	QueryTimeout                int    `json:"queryTimeout,omitempty"`
	SnapshotTransactionsEnabled bool   `json:"snapshotTransactionsEnabled,omitempty"`
	TimestampUtcEnabled         bool   `json:"timestampUtcEnabled,omitempty"`
	Timezone                    string `json:"timezone,omitempty"`
	TimeZoneBehavior            string `json:"timeZoneBehavior,omitempty"`
	ResultSetMaxRows            int    `json:"resultSetMaxRows,omitempty"`
}

type AuthCommand struct {
	Username         string     `json:"username"`
	Password         string     `json:"password"`
	UseCompression   bool       `json:"useCompression"`
	SessionID        int        `json:"sessionId,omitempty"`
	ClientName       string     `json:"clientName,omitempty"`
	DriverName       string     `json:"driverName,omitempty"`
	ClientOs         string     `json:"clientOs,omitempty"`
	ClientOsUsername string     `json:"clientOsUsername,omitempty"`
	ClientLanguage   string     `json:"clientLanguage,omitempty"`
	ClientVersion    string     `json:"clientVersion,omitempty"`
	ClientRuntime    string     `json:"clientRuntime,omitempty"`
	Attributes       Attributes `json:"attributes,omitempty"`
}

type SQLCommand struct {
	Command
	SQLText    string     `json:"sqlText"`
	Attributes Attributes `json:"attributes,omitempty"`
}
