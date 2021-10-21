package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"strings"
)

type InsertRowRequest struct {
	SchemaName string                 `json:"schemaName"`
	TableName  string                 `json:"tableName"`
	Row        map[string]interface{} `json:"row"`
}

func (request *InsertRowRequest) GetSchemaName() string {
	return "\"" + strings.ReplaceAll(request.SchemaName, "\"", "\"\"") + "\""
}

func (request *InsertRowRequest) GetTableName() string {
	return "\"" + strings.ReplaceAll(request.TableName, "\"", "\"\"") + "\""
}

func (request *InsertRowRequest) GetRow() (string, string, error) {
	var columnNames strings.Builder
	var values strings.Builder

	for columnName, value := range request.Row {
		value, err := request.getStringFromValue(value)
		if err != nil {
			return "", "", err
		}
		values.WriteString(value)
		columnNames.WriteString(fmt.Sprintf("%v,", columnName))
	}
	return strings.TrimSuffix(columnNames.String(), ","), strings.TrimSuffix(values.String(), ","), nil
}

func (request *InsertRowRequest) getStringFromValue(value interface{}) (string, error) {
	switch valueType := value.(type) {
	case bool, float32, float64, int, int8, int16, int32, int64:
		return fmt.Sprintf("%v,", value), nil
	case string:
		return "'" + strings.ReplaceAll(fmt.Sprintf("%v", value), "'", "''") + "',", nil
	default:
		return "", error_reporting_go.ExaError("E-ERA-16").
			Message("invalid row value type {{type|uq}} for value {{value|uq}} in the request").
			Parameter("type", valueType).
			Parameter("value", value)
	}
}

func (request *InsertRowRequest) Validate() error {
	if request.SchemaName == "" || request.TableName == "" || request.Row == nil {
		return error_reporting_go.ExaError("E-ERA-17").
			Message("insert row request has some missing parameters.").
			Mitigation("Please specify schema name, table name and row")
	}
	return nil
}
