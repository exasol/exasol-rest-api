package exasol_rest_api

import "github.com/tidwall/sjson"

import (
	"encoding/json"
	"fmt"
)

type GetTablesResponse struct {
	Status     string  `json:"status"`
	TablesList []Table `json:"tablesList"`
	Exception  string  `json:"exception,omitempty"`
}

type Table struct {
	TableName  string `json:"tableName"`
	SchemaName string `json:"schemaName"`
}

type GetRowsResponse struct {
	Status    string          `json:"status"`
	Rows      json.RawMessage `json:"rows,omitempty"`
	Meta      Meta            `json:"meta,omitempty"`
	Exception string          `json:"exception,omitempty"`
}

type responseData struct {
	NumResults int               `json:"numResults"`
	Results    []json.RawMessage `json:"results"`
}

type results struct {
	ResultType string    `json:"resultType"`
	ResultSet  resultSet `json:"resultSet"`
}

type resultSet struct {
	ResultSetHandle  int             `json:"resultSetHandle"`
	NumColumns       int             `json:"numColumns,omitempty"`
	NumRows          int             `json:"numRows"`
	NumRowsInMessage int             `json:"numRowsInMessage"`
	Columns          []Column        `json:"columns,omitempty"`
	Data             [][]interface{} `json:"data"`
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

type APIBaseResponse struct {
	Status    string `json:"status"`
	Exception string `json:"exception,omitempty"`
}

type webSocketsBaseResponse struct {
	Status       string          `json:"status"`
	ResponseData json.RawMessage `json:"responseData"`
	Exception    *exception      `json:"exception"`
}

type exception struct {
	Text    string `json:"text"`
	SQLCode string `json:"sqlCode"`
}

type publicKeyResponse struct {
	PublicKeyPem      string `json:"publicKeyPem"`
	PublicKeyModulus  string `json:"publicKeyModulus"`
	PublicKeyExponent string `json:"publicKeyExponent"`
}

func ConvertToGetTablesResponse(response []byte) (interface{}, error) {
	base := &webSocketsBaseResponse{}
	err := json.Unmarshal(response, base)
	if err != nil {
		return err, nil
	}

	results, err := getResults(base)
	if err != nil {
		return err, nil
	}

	convertedResponse := GetTablesResponse{
		Status: base.Status,
	}
	if base.Exception != nil {
		convertedResponse.Exception = base.Exception.SQLCode + " " + base.Exception.Text
	} else {
		convertedResponse.TablesList = []Table{}
		data := results.ResultSet.Data
		if len(data) > 0 {
			for row := range data[0] {
				convertedResponse.TablesList = append(convertedResponse.TablesList, Table{
					SchemaName: fmt.Sprintf("%v", data[0][row]),
					TableName:  fmt.Sprintf("%v", data[1][row]),
				})
			}
		}
	}
	return convertedResponse, nil
}

func ConvertToGetRowsResponse(response []byte) (interface{}, error) {
	base := &webSocketsBaseResponse{}
	err := json.Unmarshal(response, base)
	if err != nil {
		return err, nil
	}

	convertedResponse := GetRowsResponse{
		Status: base.Status,
	}
	if base.Exception != nil {
		convertedResponse.Exception = base.Exception.SQLCode + " " + base.Exception.Text
	} else {
		results, err := getResults(base)
		if err != nil {
			return err, nil
		}

		convertedResponse.Meta = Meta{Columns: results.ResultSet.Columns}
		data := results.ResultSet.Data
		if len(data) > 0 {
			rows := "["
			for rowIndex := range data[0] {
				row := ""
				for colIndex := range data {
					value := data[colIndex][rowIndex]
					row, _ = sjson.Set(row, convertedResponse.Meta.Columns[colIndex].Name, value)
				}
				rows += row
				if rowIndex < len(data[0])-1 {
					rows += ","
				}
			}
			rows += "]"
			convertedResponse.Rows = json.RawMessage(rows)
		}
	}
	return convertedResponse, nil
}

func getResults(base *webSocketsBaseResponse) (*results, error) {
	responseData := &responseData{}
	err := json.Unmarshal(base.ResponseData, responseData)
	if err != nil {
		return nil, err
	}

	results := &results{}
	err = json.Unmarshal(responseData.Results[0], results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func ConvertToBaseResponse(response []byte) (interface{}, error) {
	base := &webSocketsBaseResponse{}
	err := json.Unmarshal(response, base)
	if err != nil {
		return err, nil
	}

	convertedResponse := APIBaseResponse{
		Status: base.Status,
	}
	if base.Exception != nil {
		convertedResponse.Exception = base.Exception.SQLCode + " " + base.Exception.Text
	}
	return convertedResponse, nil
}
