package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"strings"
)

// InsertRowRequest maps an InsertRow JSON request to a struct.
type InsertRowRequest struct {
	SchemaName string  `json:"schemaName"`
	TableName  string  `json:"tableName"`
	Row        []Value `json:"row"`
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

	for index, value := range request.Row {
		renderedValue, err := value.getValue()
		if err != nil {
			return "", "", err
		}
		values.WriteString(renderedValue)
		columnNames.WriteString(fmt.Sprintf("%v", value.getColumnName()))
		if index < len(request.Row)-1 {
			values.WriteString(",")
			columnNames.WriteString(",")
		}
	}
	return columnNames.String(), values.String(), nil
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
