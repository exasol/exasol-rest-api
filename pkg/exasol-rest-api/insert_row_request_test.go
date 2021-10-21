package exasol_rest_api_test

import (
	"github.com/stretchr/testify/suite"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"testing"
)

type InsertRowRequestSuite struct {
	suite.Suite
}

func TestInsertRowRequestSuite(t *testing.T) {
	suite.Run(t, new(InsertRowRequestSuite))
}

func (suite *InsertRowRequestSuite) TestGetSchemaName() {
	request := exasol_rest_api.InsertRowRequest{
		SchemaName: "MY_SCHEMA",
	}
	suite.Equal("\"MY_SCHEMA\"", request.GetSchemaName())
}

func (suite *InsertRowRequestSuite) TestGetTableName() {
	request := exasol_rest_api.InsertRowRequest{
		TableName: "MY_TABLE",
	}
	suite.Equal("\"MY_TABLE\"", request.GetTableName())
}

func (suite *InsertRowRequestSuite) TestGetRow() {
	request := exasol_rest_api.InsertRowRequest{
		Row: map[string]interface{}{
			"c1": "Exa'sol",
			"c2": 3,
			"c3": 123.456,
			"c4": false,
			"c5": "3 12:50:10.123",
		},
	}
	columns, values, _ := request.GetRow()
	suite.Contains(columns, "c1")
	suite.Contains(columns, "c2")
	suite.Contains(columns, "c3")
	suite.Contains(columns, "c4")
	suite.Contains(columns, "c5")
	suite.Contains(values, "'Exa''sol'")
	suite.Contains(values, "3")
	suite.Contains(values, "123.456")
	suite.Contains(values, "false")
	suite.Contains(values, "'3 12:50:10.123'")
}

func (suite *InsertRowRequestSuite) TestGetRowWithInvalidInterface() {
	request := exasol_rest_api.InsertRowRequest{
		Row: map[string]interface{}{
			"c1": nil,
		},
	}
	columns, values, err := request.GetRow()
	suite.EqualError(err,
		"E-ERA-16: invalid row value type <nil> for value <nil> in the request")
	suite.Equal(columns, "")
	suite.Equal(values, "")
}

func (suite *InsertRowRequestSuite) TestValidateSuccess() {
	request := exasol_rest_api.InsertRowRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		Row:        map[string]interface{}{"key": "value"},
	}
	suite.NoError(request.Validate())
}

func (suite *InsertRowRequestSuite) TestValidateWithoutSchemaName() {
	request := exasol_rest_api.InsertRowRequest{
		TableName: "MY_TABLE",
		Row:       map[string]interface{}{"key": "value"},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-17: insert row request has some missing parameters. "+
			"Please specify schema name, table name and row")
}

func (suite *InsertRowRequestSuite) TestValidateWithoutTableName() {
	request := exasol_rest_api.InsertRowRequest{
		SchemaName: "MY_SCHEMA",
		Row:        map[string]interface{}{"key": "value"},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-17: insert row request has some missing parameters. "+
			"Please specify schema name, table name and row")
}

func (suite *InsertRowRequestSuite) TestValidateWithoutRow() {
	request := exasol_rest_api.InsertRowRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
	}
	suite.EqualError(request.Validate(),
		"E-ERA-17: insert row request has some missing parameters. "+
			"Please specify schema name, table name and row")
}
