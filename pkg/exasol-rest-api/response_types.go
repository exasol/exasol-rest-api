package exasol_rest_api

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
	Columns          []column        `json:"columns,omitempty"`
	Data             [][]interface{} `json:"data"`
}

type column struct {
	Name     string   `json:"name"`
	DataType dataType `json:"dataType"`
}

type dataType struct {
	Type              string  `json:"type"`
	Precision         *int64  `json:"precision,omitempty"`
	Scale             *int64  `json:"scale,omitempty"`
	Size              *int64  `json:"size,omitempty"`
	CharacterSet      *string `json:"characterSet,omitempty"`
	WithLocalTimeZone *bool   `json:"withLocalTimeZone,omitempty"`
	Fraction          *int    `json:"fraction,omitempty"`
	SRID              *int    `json:"srid,omitempty"`
}

type baseResponse struct {
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
	base := &baseResponse{}
	err := json.Unmarshal(response, base)
	if err != nil {
		return err, nil
	}

	responseData := &responseData{}
	err = json.Unmarshal(base.ResponseData, responseData)
	if err != nil {
		return err, nil
	}

	results := &results{}
	err = json.Unmarshal(responseData.Results[0], results)
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
		if data != nil && len(data) > 0 {
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
