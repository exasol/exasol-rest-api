package exasol_rest_api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/sjson"
)

var statusOk = "ok"

// [impl->dsn~get-tables-response-body~1]
type GetTablesResponse struct {
	Status     string  `json:"status"`
	TablesList []Table `json:"tablesList"`
	Exception  string  `json:"exception,omitempty"`
}

type Table struct {
	TableName  string `json:"tableName"`
	SchemaName string `json:"schemaName"`
}

// [impl->dsn~execute-query-response-body~1]
// [impl->dsn~get-rows-response-body~1]
type GetRowsResponse struct {
	Status    string          `json:"status"`
	Rows      json.RawMessage `json:"rows,omitempty"`
	Meta      Meta            `json:"meta,omitempty"`
	Exception string          `json:"exception,omitempty"`
}

type Meta struct {
	Columns []Column `json:"columns,omitempty"`
}

type Column struct {
	Name     string   `json:"name"`
	DataType DataType `json:"dataType"`
}

type DataType struct {
	Type              string `json:"type"`
	Precision         int64  `json:"precision,omitempty"`
	Scale             int64  `json:"scale,omitempty"`
	Size              int64  `json:"size,omitempty"`
	CharacterSet      string `json:"characterSet,omitempty"`
	WithLocalTimeZone bool   `json:"withLocalTimeZone,omitempty"`
	Fraction          int    `json:"fraction,omitempty"`
	SRID              int    `json:"srid,omitempty"`
}

// [impl->dsn~insert-row-response-body~1]
// [impl->dsn~delete-rows-response-body~1]
// [impl->dsn~update-rows-response-body~1]
// [impl->dsn~execute-statement-response-body~1]
type APIBaseResponse struct {
	Status    string `json:"status"`
	Exception string `json:"exception,omitempty"`
}

// [impl->dsn~get-tables-response-body~1]
func ConvertToGetTablesResponse(rows *sql.Rows) (interface{}, error) {
	convertedResponse := GetTablesResponse{
		Status:     statusOk,
		TablesList: []Table{},
	}
	for rows.Next() {
		var table Table
		err := rows.Scan(&table.SchemaName, &table.TableName)
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}
		convertedResponse.TablesList = append(convertedResponse.TablesList, table)
	}
	return convertedResponse, nil
}

// [impl->dsn~execute-query-response-body~1]
func ConvertToGetRowsResponse(rows *sql.Rows) (interface{}, error) {
	columns, err := extractColumns(rows)
	if err != nil {
		return nil, err
	}
	rowsJson, err := buildRowsString(rows)
	if err != nil {
		return nil, err
	}
	convertedResponse := GetRowsResponse{
		Status: statusOk,
		Meta:   Meta{Columns: columns},
		Rows:   json.RawMessage(rowsJson),
	}

	return convertedResponse, nil
}

func extractColumns(rows *sql.Rows) ([]Column, error) {
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	names, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if len(types) != len(names) {
		return nil, fmt.Errorf("inconsistent row metadata: %d types and %d names", len(types), len(names))
	}
	columns := []Column{}
	for i := range names {
		columns = append(columns, createColumn(names[i], types[i]))
	}
	return columns, nil
}

func createColumn(colName string, colType *sql.ColumnType) Column {
	fmt.Printf("Col %s: type %v\n", colName, colType)
	precision, scale, _ := colType.DecimalSize()
	length, _ := colType.Length()
	return Column{
		Name: colName,
		DataType: DataType{
			Type:              colType.DatabaseTypeName(),
			Precision:         precision,
			Scale:             scale,
			Size:              length,
			CharacterSet:      "",
			WithLocalTimeZone: false,
			Fraction:          0,
			SRID:              0,
		},
	}
}

func buildRowsString(sqlRows *sql.Rows) (string, error) {
	rows := "["
	types, err := sqlRows.ColumnTypes()
	if err != nil {
		return "", err
	}
	names, err := sqlRows.Columns()
	if err != nil {
		return "", err
	}

	dest := []any{}
	for _, colType := range types {
		dest = append(dest, destForType(colType))
	}

	for sqlRows.Next() {
		err = sqlRows.Scan(dest...)
		if err != nil {
			return "", err
		}
		row := ""
		for colIndex := range dest {
			value := dest[colIndex]
			row, _ = sjson.Set(row, names[colIndex], &value)
			fmt.Printf("Col %d: %s = %s / %v ---- Row = %s\n", colIndex, names[colIndex], value, &value, row)
		}
		rows += row
		rows += ","
	}
	rows = strings.TrimRight(rows, ",")
	rows += "]"
	return rows, nil
}

func destForType(colType *sql.ColumnType) any {
	t := colType.ScanType()
	fmt.Printf("Col %s: %v\n", colType.Name(), t)
	//dest := ""
	return new(any)
}
