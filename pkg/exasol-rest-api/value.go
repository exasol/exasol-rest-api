package exasol_rest_api

// Value represents a single cell in a table.
type Value struct {
	ColumnName string      `json:"columnName"`
	Value      interface{} `json:"value"`
}

func (value *Value) getColumnName() string {
	return toExasolIdentifier(value.ColumnName)
}

func (value *Value) render() (string, error) {
	return renderLiteral(value.Value)
}

func (value *Value) validate() bool {
	return value.ColumnName != "" && value.Value != nil
}
