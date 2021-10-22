package exasol_rest_api_test

import (
	"github.com/stretchr/testify/suite"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"testing"
)

type UpdateRowsRequestSuite struct {
	suite.Suite
}

func TestUpdateRowsRequestSuite(t *testing.T) {
	suite.Run(t, new(UpdateRowsRequestSuite))
}

func (suite *UpdateRowsRequestSuite) TestGetSchemaName() {
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: "MY_SCHEMA",
	}
	suite.Equal("\"MY_SCHEMA\"", request.GetSchemaName())
}

func (suite *UpdateRowsRequestSuite) TestGetTableName() {
	request := exasol_rest_api.UpdateRowsRequest{
		TableName: "MY_TABLE",
	}
	suite.Equal("\"MY_TABLE\"", request.GetTableName())
}

func (suite *UpdateRowsRequestSuite) TestGetValuesToUpdate() {
	request := exasol_rest_api.UpdateRowsRequest{
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "c1", Value: "Exa'sol"},
			{ColumnName: "c2", Value: 3},
			{ColumnName: "c3", Value: 123.456},
			{ColumnName: "c4", Value: false},
			{ColumnName: "c5", Value: "3 12:50:10.123"},
		},
	}
	values, _ := request.GetValuesToUpdate()
	suite.Equal("\"c1\"='Exa''sol',\"c2\"=3,\"c3\"=123.456,\"c4\"=false,\"c5\"='3 12:50:10.123'", values)
}

func (suite *UpdateRowsRequestSuite) TestGetRowWithInvalidInterface() {
	request := exasol_rest_api.UpdateRowsRequest{
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "c1", Value: nil},
		},
	}
	values, err := request.GetValuesToUpdate()
	suite.EqualError(err,
		"E-ERA-16: invalid exasol literal type <nil> for value <nil> in the request")
	suite.Equal("", values)
}

func (suite *UpdateRowsRequestSuite) TestValidateSuccess() {
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.NoError(request.Validate())
}

func (suite *UpdateRowsRequestSuite) TestValidateWithoutSchemaName() {
	request := exasol_rest_api.UpdateRowsRequest{
		TableName: "MY_TABLE",
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-20: update rows request has some missing parameters. "+
			"Please specify schema name, table name, values to update and condition")
}

func (suite *UpdateRowsRequestSuite) TestValidateWithoutTableName() {
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: "MY_SCHEMA",
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "key", Value: "value"},
		},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-20: update rows request has some missing parameters. "+
			"Please specify schema name, table name, values to update and condition")
}

func (suite *UpdateRowsRequestSuite) TestValidateWithoutValuesToUpdate() {
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-20: update rows request has some missing parameters. "+
			"Please specify schema name, table name, values to update and condition")
}

func (suite *UpdateRowsRequestSuite) TestValidateWithRowsToUpdateEmpty() {
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName:     "MY_SCHEMA",
		TableName:      "MY_TABLE",
		ValuesToUpdate: []exasol_rest_api.Value{},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-20: update rows request has some missing parameters. "+
			"Please specify schema name, table name, values to update and condition")
}

func (suite *UpdateRowsRequestSuite) TestValidateWithRowsToUpdateWithoutColumnName() {
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		ValuesToUpdate: []exasol_rest_api.Value{
			{Value: "value"},
		},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-20: update rows request has some missing parameters. "+
			"Please specify schema name, table name, values to update and condition")
}

func (suite *UpdateRowsRequestSuite) TestValidateWithRowsToUpdateWithoutValue() {
	request := exasol_rest_api.UpdateRowsRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		ValuesToUpdate: []exasol_rest_api.Value{
			{ColumnName: "key"},
		},
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-20: update rows request has some missing parameters. "+
			"Please specify schema name, table name, values to update and condition")
}
