package exasol_rest_api

import "fmt"

// Condition represents a simple SQL WHERE condition.
type Condition struct {
	CellValue           Value  `json:"value"`
	ComparisonPredicate string `json:"comparisonPredicate"`
}

func (whereCondition *Condition) getColumnName() string {
	return whereCondition.CellValue.getColumnName()
}

func (whereCondition *Condition) getComparisonPredicate() string {
	predicate := whereCondition.ComparisonPredicate
	if predicate == "=" || predicate == "!=" ||
		predicate == ">" || predicate == "<" ||
		predicate == ">=" || predicate == "<=" {
		return predicate
	} else {
		return "="
	}
}

func (whereCondition *Condition) getValue() (string, error) {
	return whereCondition.CellValue.render()
}

func (whereCondition *Condition) validate() bool {
	return whereCondition.CellValue.validate() && (whereCondition.ComparisonPredicate == "" ||
		whereCondition.ComparisonPredicate == "=" || whereCondition.ComparisonPredicate == "!=" ||
		whereCondition.ComparisonPredicate == ">" || whereCondition.ComparisonPredicate == "<" ||
		whereCondition.ComparisonPredicate == ">=" || whereCondition.ComparisonPredicate == "<=")
}

func (whereCondition *Condition) render() (string, error) {
	columnName := whereCondition.getColumnName()
	comparisonPredicate := whereCondition.getComparisonPredicate()
	columnValue, err := whereCondition.getValue()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v %v %v", columnName, comparisonPredicate, columnValue), nil
}
