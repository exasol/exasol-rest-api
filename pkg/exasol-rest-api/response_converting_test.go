package exasol_rest_api_test

import (
	"github.com/stretchr/testify/suite"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"testing"
)

type ResponseConvertingSuite struct {
	suite.Suite
}

func TestResponseConvertingSuite(t *testing.T) {
	suite.Run(t, new(ResponseConvertingSuite))
}

func (suite *ResponseConvertingSuite) TestConvertGetTablesResponse() {
	converted, err := exasol_rest_api.ConvertToGetTablesResponse([]byte("{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":2,\"numRowsInMessage\":2,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}}],\"data\":[[\"TEST_SCHEMA_DELETE_ROWS_1\",\"TEST_SCHEMA_DELETE_ROWS_1\"],[\"TEST_TABLE\",\"TEST_TABLE_1\"]]}}],\"numResults\":1}}"))
	expected := exasol_rest_api.GetTablesResponse{
		Status: "ok",
		TablesList: []exasol_rest_api.Table{
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
	converted, err := exasol_rest_api.ConvertToGetTablesResponse([]byte("{\"status\":\"ok\",\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":0,\"numRowsInMessage\":0,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}}]}}],\"numResults\":1}}"))
	expected := exasol_rest_api.GetTablesResponse{
		Status:     "ok",
		TablesList: []exasol_rest_api.Table{},
	}
	suite.Equal(expected, converted)
	suite.NoError(err)
}

func (suite *ResponseConvertingSuite) TestConvertGetTablesResponseWithError() {
	converted, err := exasol_rest_api.ConvertToGetTablesResponse([]byte("{\"status\":\"error\",\"exception\":{\"sqlCode\":\"1\",\"text\":\"message\"},\"responseData\":{\"results\":[{\"resultType\":\"resultSet\",\"resultSet\":{\"numColumns\":2,\"numRows\":0,\"numRowsInMessage\":0,\"columns\":[{\"name\":\"TABLE_SCHEMA\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}},{\"name\":\"TABLE_NAME\",\"dataType\":{\"type\":\"VARCHAR\",\"size\":128,\"characterSet\":\"UTF8\"}}]}}],\"numResults\":1}}"))
	expected := exasol_rest_api.GetTablesResponse{
		Status:     "error",
		TablesList: nil,
		Exception:  "1 message",
	}
	suite.Equal(expected, converted)
	suite.NoError(err)
}
