package exasol_rest_api

import (
	error_reporting_go "github.com/exasol/error-reporting-go"
)

// RowsRequest maps DeleteRows and GetRows requests to a struct.
// [impl->dsn~delete-rows-request-body~1]
// [impl->dsn~get-rows-request-parameters~1]
type RowsRequest struct {
	SchemaName     string    `json:"schemaName"`
	TableName      string    `json:"tableName"`
	WhereCondition Condition `json:"condition"`
}

// GetSchemaName returns a schema name.
func (request *RowsRequest) GetSchemaName() string {
	return toExasolIdentifier(request.SchemaName)
}

// GetTableName returns a table name.
func (request *RowsRequest) GetTableName() string {
	return toExasolIdentifier(request.TableName)
}

// GetCondition return a rendered condition.
func (request *RowsRequest) GetCondition() (string, error) {
	return request.WhereCondition.render()
}

// Validate validates the request.
func (request *RowsRequest) Validate() error {
	if request.SchemaName == "" || request.TableName == "" || !request.WhereCondition.validate() {
		return error_reporting_go.ExaError("E-ERA-19").
			Message("request has some missing parameters.").
			Mitigation("Please specify schema name, table name and condition: column name, value")
	}
	return nil
}