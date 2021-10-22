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
		Row: []exasol_rest_api.Value{
			{ColumnName: "c1", Value: "Exa'sol"},
			{ColumnName: "c2", Value: 3},
			{ColumnName: "c3", Value: 123.456},
			{ColumnName: "c4", Value: false},
			{ColumnName: "c5", Value: "3 12:50:10.123"},
		},
	}
	columns, values, _ := request.GetRow()
	suite.Equal("\"c1\",\"c2\",\"c3\",\"c4\",\"c5\"", columns)
	suite.Equal("'Exa''sol',3,123.456,false,'3 12:50:10.123'", values)
}

func (suite *InsertRowRequestSuite) TestGetRowWithInvalidInterface() {
	request := exasol_rest_api.InsertRowRequest{
		Row: []exasol_rest_api.Value{
			{ColumnName: "c1", Value: nil},
		},
	}
	columns, values, err := request.GetRow()
	suite.EqualError(err,
		"E-ERA-16: invalid exasol literal type <nil> for value <nil> in the request")
	suite.Equal("", columns)
	suite.Equal("", values)
}

func (suite *InsertRowRequestSuite) TestValidateSuccess() {
	request := exasol_rest_api.InsertRowRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		Row: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
	}
	suite.NoError(request.Validate())
}

func (suite *InsertRowRequestSuite) TestValidateWithoutSchemaName() {
	request := exasol_rest_api.InsertRowRequest{
		TableName: "MY_TABLE",
		Row: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-17: insert row request has some missing parameters. "+
			"Please specify schema name, table name and row")
}

func (suite *InsertRowRequestSuite) TestValidateWithoutTableName() {
	request := exasol_rest_api.InsertRowRequest{
		SchemaName: "MY_SCHEMA",
		Row: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
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
