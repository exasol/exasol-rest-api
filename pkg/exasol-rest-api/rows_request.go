package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
)

// RowsRequest maps a DeleteRows, UpdateRows and GetRows JSON requests to a struct.
type RowsRequest struct {
	SchemaName     string    `json:"schemaName"`
	TableName      string    `json:"tableName"`
	WhereCondition Condition `json:"condition"`
}

// Condition represents a simple SQL WHERE condition.
type Condition struct {
	ColumnName          string      `json:"columnName"`
	ColumnValue         interface{} `json:"columnValue"`
	ComparisonPredicate string      `json:"comparisonPredicate"`
}

func (whereCondition *Condition) getColumnName() string {
	return ToExasolIdentifier(whereCondition.ColumnName)
}

func (whereCondition *Condition) getComparisonPredicate() (string, error) {
	predicate := whereCondition.ComparisonPredicate
	if predicate == "" {
		return "=", nil
	} else if predicate == "=" || predicate == "!=" ||
		predicate == ">" || predicate == "<" ||
		predicate == ">=" || predicate == "<=" {
		return predicate, nil
	}
	return "", error_reporting_go.ExaError("E-ERA-18").
		Message("invalid predicate value: {{predicate}}.").
		Parameter("predicate", predicate).
		Mitigation("Please use one of the following values: =, !=, <, >, <=, >=")
}

func (whereCondition *Condition) getColumnValue() (string, error) {
	return ToExasolLiteral(whereCondition.ColumnValue)
}

// GetSchemaName returns a schema name.
func (request *RowsRequest) GetSchemaName() string {
	return ToExasolIdentifier(request.SchemaName)
}

// GetTableName returns a table name.
func (request *RowsRequest) GetTableName() string {
	return ToExasolIdentifier(request.TableName)
}

// GetCondition return a rendered condition.
func (request *RowsRequest) GetCondition() (string, error) {
	columnName := request.WhereCondition.getColumnName()
	comparisonPredicate, err := request.WhereCondition.getComparisonPredicate()
	if err != nil {
		return "", err
	}
	columnValue, err := request.WhereCondition.getColumnValue()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v %v %v", columnName, comparisonPredicate, columnValue), nil
}

// Validate validates the request.
func (request *RowsRequest) Validate() error {
	if request.SchemaName == "" || request.TableName == "" ||
		request.WhereCondition.ColumnName == "" || request.WhereCondition.ColumnValue == nil {
		return error_reporting_go.ExaError("E-ERA-19").
			Message("request has some missing parameters.").
			Mitigation("Please specify schema name, table name and condition: column name, value")
	}
	return nil
}
