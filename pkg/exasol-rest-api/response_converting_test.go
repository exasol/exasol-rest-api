package exasol_rest_api_test

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	. "main/pkg/exasol-rest-api"
	"testing"
)

type ResponseConvertingSuite struct {
	suite.Suite
}

// [utest->dsn~get-tables-response-body~1]
// [utest->dsn~get-rows-response-body~1]
// [utest->dsn~execute-query-response-body~1]
func TestResponseConvertingSuite(t *testing.T) {
	suite.Run(t, new(ResponseConvertingSuite))
}

func (suite *ResponseConvertingSuite) TestConvertGetTablesResponse() {
	converted, err := ConvertToGetTablesResponse([]byte("{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":2,\"numRowsInMessage\":2,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}}],\"data\":[[\"TEST_SCHEMA_DELETE_ROWS_1\",\"TEST_SCHEMA_DELETE_ROWS_1\"],[\"TEST_TABLE\",\"TEST_TABLE_1\"]]}}],\"numResults\":1}}"))
	expected := GetTablesResponse{
		Status: "ok",
		TablesList: []Table{
			{
				TableName:  "TEST_TABLE",
				SchemaName: "TEST_SCHEMA_DELETE_ROWS_1",
			},
			{
				TableName:  "TEST_TABLE_1",
				SchemaName: "TEST_SCHEMA_DELETE_ROWS_1",
			},
		},
	}
	suite.Equal(expected, converted)
	suite.NoError(err)
}

func (suite *ResponseConvertingSuite) TestConvertGetTablesResponseWithZeroTables() {
	converted, err := ConvertToGetTablesResponse([]byte("{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":0,\"numRowsInMessage\":0,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}}]}}],\"numResults\":1}}"))
	expected := GetTablesResponse{
		Status:     "ok",
		TablesList: []Table{},
	}
	suite.Equal(expected, converted)
	suite.NoError(err)
}

func (suite *ResponseConvertingSuite) TestConvertGetTablesResponseWithError() {
	converted, err := ConvertToGetTablesResponse([]byte("{\"status\":\"error\",\"exception\":{\"sqlCode\":\"1\",\"text\":\"message\"},\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":0,\"numRowsInMessage\":0,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}}]}}],\"numResults\":1}}"))
	expected := GetTablesResponse{
		Status:     "error",
		TablesList: nil,
		Exception:  "1 message",
	}
	suite.Equal(expected, converted)
	suite.NoError(err)
}

func (suite *ResponseConvertingSuite) TestConvertGetRowsResponse() {
	converted, err := ConvertToGetRowsResponse([]byte("{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":2,\"numRowsInMessage\":2,\"columns\":[{\"name\":\"X\",\"dataType\":{\"type\":\"DECIMAL\",\"precision\":18,\"scale\":0}},{\"name\":\"Y\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":100,\"characterSet\":\"UTF8\"}}],\"data\":[[15,10],[\"test\",\"foo\"]]}}],\"numResults\":1}}"))
	expected := GetRowsResponse{
		Status: "ok",
		Meta: Meta{
			Columns: []Column{
				{Name: "X", DataType: DataType{Type: "DECIMAL", Precision: int64(18), Scale: int64(0)}},
				{Name: "Y", DataType: DataType{Type: "VARCHAR", Size: int64(100), CharacterSet: "UTF8"}},
			},
		},
		Rows: json.RawMessage("[{\"X\":15,\"Y\":\"test\"},{\"X\":10,\"Y\":\"foo\"}]"),
	}
	suite.Equal(expected, converted)
	suite.NoError(err)
}
