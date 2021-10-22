package exasol_rest_api

import (
	error_reporting_go "github.com/exasol/error-reporting-go"
)

// DeleteRowsRequest maps a DeleteRows JSON request to a struct.
type DeleteRowsRequest struct {
	SchemaName     string    `json:"schemaName"`
	TableName      string    `json:"tableName"`
	WhereCondition Condition `json:"condition"`
}

// GetSchemaName returns a schema name.
func (request *DeleteRowsRequest) GetSchemaName() string {
	return ToExasolIdentifier(request.SchemaName)
}

// GetTableName returns a table name.
func (request *DeleteRowsRequest) GetTableName() string {
	return ToExasolIdentifier(request.TableName)
}

// GetCondition return a rendered condition.
func (request *DeleteRowsRequest) GetCondition() (string, error) {
	return renderCondition(request.WhereCondition)
}

// Validate validates the request.
func (request *DeleteRowsRequest) Validate() error {
	if request.SchemaName == "" || request.TableName == "" || !request.WhereCondition.validate() {
		return error_reporting_go.ExaError("E-ERA-19").
			Message("request has some missing parameters.").
			Mitigation("Please specify schema name, table name and condition: column name, value")
	}
	return nil
}
