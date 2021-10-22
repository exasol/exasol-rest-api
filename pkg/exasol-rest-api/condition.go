package exasol_rest_api

import error_reporting_go "github.com/exasol/error-reporting-go"

// Condition represents a simple SQL WHERE condition.
type Condition struct {
	CellValue           Value  `json:"value"`
	ComparisonPredicate string `json:"comparisonPredicate"`
}

func (whereCondition *Condition) getColumnName() string {
	return whereCondition.CellValue.getColumnName()
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

func (whereCondition *Condition) getValue() (string, error) {
	return whereCondition.CellValue.getValue()
}

func (whereCondition *Condition) validate() bool {
	return whereCondition.CellValue.validate()
}
