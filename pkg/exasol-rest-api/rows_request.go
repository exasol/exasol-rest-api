package exasol_rest_api

import (
	exaerror "github.com/exasol/error-reporting-go"
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
		return createValidationError().
			Mitigation("Please specify schema name, table name and condition: column name, value")

	}
	return nil
}

func createValidationError() *exaerror.ExaError {
	return exaerror.New("E-ERA-19").Message("request has some missing parameters.")
}

// Validate validates the request when the condition is optional.
func (request *RowsRequest) ValidateWithOptionalCondition() error {
	if request.SchemaName == "" || request.TableName == "" {
		return createValidationError().
			Mitigation("Please specify schema name and table name")
	}
	return nil
}

func (request *RowsRequest) HasWhereClause() bool {
	return request.WhereCondition.validate()
}
