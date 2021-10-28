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
	return toExasolIdentifier(request.SchemaName)
}

// GetTableName returns a table name.
func (request *DeleteRowsRequest) GetTableName() string {
	return toExasolIdentifier(request.TableName)
}

// GetCondition return a rendered condition.
func (request *DeleteRowsRequest) GetCondition() (string, error) {
	return request.WhereCondition.render()
}

// Validate validates the request.
func (request *DeleteRowsRequest) Validate() error {
	if request.SchemaName == "" || request.TableName == "" || !request.WhereCondition.validate() {
		return error_reporting_go.ExaError("E-ERA-19").
			Message("delete rows request has some missing parameters.").
			Mitigation("Please specify schema name, table name and condition: column name, value")
	}
	return nil
}
