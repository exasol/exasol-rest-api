package exasol_rest_api

import (
	"fmt"
	error_reporting_go "github.com/exasol/error-reporting-go"
	"strings"
)

func renderCondition(condition Condition) (string, error) {
	columnName := condition.getColumnName()
	comparisonPredicate := condition.getComparisonPredicate()
	columnValue, err := condition.getValue()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v %v %v", columnName, comparisonPredicate, columnValue), nil
}

func renderLiteral(value interface{}) (string, error) {
	switch valueType := value.(type) {
	case bool, float32, float64, int, int8, int16, int32, int64:
		return fmt.Sprintf("%v", value), nil
	case string:
		return "'" + strings.ReplaceAll(fmt.Sprintf("%v", value), "'", "''") + "'", nil
	default:
		return "", error_reporting_go.ExaError("E-ERA-16").
			Message("invalid exasol literal type {{type|uq}} for value {{value|uq}} in the request").
			Parameter("type", valueType).
			Parameter("value", value)
	}
}
