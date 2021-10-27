package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"strings"
)

// UpdateRowsRequest maps an UpdateRows JSON request to a struct.
type UpdateRowsRequest struct {
	SchemaName     string    `json:"schemaName"`
	TableName      string    `json:"tableName"`
	ValuesToUpdate []Value   `json:"row"`
	WhereCondition Condition `json:"condition"`
}

// GetSchemaName returns a schema name.
func (request *UpdateRowsRequest) GetSchemaName() string {
	return toExasolIdentifier(request.SchemaName)
}

// GetTableName returns a table name.
func (request *UpdateRowsRequest) GetTableName() string {
	return toExasolIdentifier(request.TableName)
}

// GetCondition returns a rendered condition.
func (request *UpdateRowsRequest) GetCondition() (string, error) {
	return request.WhereCondition.render()
}

// GetCondition returns a rendered condition.
func (request *UpdateRowsRequest) GetValuesToUpdate() (string, error) {
	var valuesToUpdate strings.Builder

	for index, value := range request.ValuesToUpdate {
		renderedValue, err := value.render()
		if err != nil {
			return "", err
		}
		valuesToUpdate.WriteString(fmt.Sprintf("%v=%v", value.getColumnName(), renderedValue))
		if index < len(request.ValuesToUpdate)-1 {
			valuesToUpdate.WriteString(",")
		}
	}
	return valuesToUpdate.String(), nil
}

// Validate validates the request.
func (request *UpdateRowsRequest) Validate() error {
	valuesValidation := true
	for _, value := range request.ValuesToUpdate {
		if !value.validate() {
			valuesValidation = false
			break
		}
	}

	if request.SchemaName == "" || request.TableName == "" || !request.WhereCondition.validate() ||
		request.ValuesToUpdate == nil || len(request.ValuesToUpdate) == 0 || !valuesValidation {
		return error_reporting_go.ExaError("E-ERA-20").
			Message("update rows request has some missing parameters.").
			Mitigation("Please specify schema name, table name, values to update and condition")
	}
	return nil
}
