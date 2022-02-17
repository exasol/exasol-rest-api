package exasol_rest_api

import (
	"fmt"
	"strings"

	exaerror "github.com/exasol/error-reporting-go"
)

// InsertRowRequest maps an InsertRow JSON request to a struct.
// [impl->dsn~insert-row-request-body~1]
type InsertRowRequest struct {
	SchemaName string  `json:"schemaName"`
	TableName  string  `json:"tableName"`
	Row        []Value `json:"row"`
}

// GetSchemaName returns a schema name.
func (request *InsertRowRequest) GetSchemaName() string {
	return toExasolIdentifier(request.SchemaName)
}

// GetTableName returns a table name.
func (request *InsertRowRequest) GetTableName() string {
	return toExasolIdentifier(request.TableName)
}

// GetRow returns column names and values of the row.
func (request *InsertRowRequest) GetRow() (string, string, error) {
	var columnNames strings.Builder
	var values strings.Builder

	for index, value := range request.Row {
		renderedValue, err := value.render()
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
	if request.SchemaName == "" || request.TableName == "" || request.Row == nil || len(request.Row) == 0 {
		return exaerror.New("E-ERA-17").
			Message("insert row request has some missing parameters.").
			Mitigation("Please specify schema name, table name and row")
	}
	return nil
}
