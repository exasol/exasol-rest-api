package exasol_rest_api

import exaerror "github.com/exasol/error-reporting-go"

// ExecuteStatementRequest maps an ExecuteStatement JSON request to a struct.
// [impl->dsn~execute-statement-request-body~1]
type ExecuteStatementRequest struct {
	Statement string `json:"sqlStatement"`
}

// GetStatement returns a statement.
func (request *ExecuteStatementRequest) GetStatement() string {
	return request.Statement
}

// Validate validates the request.
func (request *ExecuteStatementRequest) Validate() error {
	if request.Statement == "" {
		return exaerror.New("E-ERA-29").
			Message("execute statement request has a missing statement.").
			Mitigation("Please add a statement to the request body")
	}
	return nil
}
