package exasol_rest_api

import error_reporting_go "github.com/exasol/error-reporting-go"

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
	return renderLiteral(whereCondition.ColumnValue)
}

func (whereCondition *Condition) validate() bool {
	return whereCondition.ColumnName != "" && whereCondition.ColumnValue != nil
}
