package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"strings"
)

// InsertRowRequest maps an InsertRow JSON request to a struct.
type InsertRowRequest struct {
	SchemaName string                 `json:"schemaName"`
	TableName  string                 `json:"tableName"`
	Row        map[string]interface{} `json:"row"`
}

// GetSchemaName returns a schema name.
func (request *InsertRowRequest) GetSchemaName() string {
	return ToExasolIdentifier(request.SchemaName)
}

// GetTableName returns a table name.
func (request *InsertRowRequest) GetTableName() string {
	return ToExasolIdentifier(request.TableName)
}

//GetRow returns columns names and values of the row.
func (request *InsertRowRequest) GetRow() (string, string, error) {
	var columnNames strings.Builder
	var values strings.Builder

	for columnName, value := range request.Row {
		value, err := ToExasolLiteral(value)
		if err != nil {
			return "", "", err
		}
		values.WriteString(value)
		values.WriteString(",")
		columnNames.WriteString(fmt.Sprintf("%v,", columnName))
	}
	return strings.TrimSuffix(columnNames.String(), ","), strings.TrimSuffix(values.String(), ","), nil
}

// Validate validates the request.
func (request *InsertRowRequest) Validate() error {
	if request.SchemaName == "" || request.TableName == "" || request.Row == nil {
		return error_reporting_go.ExaError("E-ERA-17").
			Message("insert row request has some missing parameters.").
			Mitigation("Please specify schema name, table name and row")
	}
	return nil
}
