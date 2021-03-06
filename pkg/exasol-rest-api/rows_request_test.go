package exasol_rest_api_test

import (
	"github.com/stretchr/testify/suite"
	exasol_rest_api "main/pkg/exasol-rest-api"
	"testing"
)

type RowsRequestSuite struct {
	suite.Suite
}

// [utest->dsn~delete-rows-request-body~1]
// [utest->dsn~get-rows-request-parameters~1]
func TestRowsRequestSuite(t *testing.T) {
	suite.Run(t, new(RowsRequestSuite))
}

func (suite *RowsRequestSuite) TestGetSchemaName() {
	request := exasol_rest_api.RowsRequest{
		SchemaName: "MY_SCHEMA",
	}
	suite.Equal("\"MY_SCHEMA\"", request.GetSchemaName())
}

func (suite *RowsRequestSuite) TestGetTableName() {
	request := exasol_rest_api.RowsRequest{
		TableName: "MY_TABLE",
	}
	suite.Equal("\"MY_TABLE\"", request.GetTableName())
}

func (suite *RowsRequestSuite) TestGetCondition() {
	request := exasol_rest_api.RowsRequest{
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      100,
			},
			ComparisonPredicate: "<",
		},
	}
	condition, err := request.GetCondition()
	suite.Equal("\"MY_COLUMN\" < 100", condition)
	suite.NoError(err)
}

func (suite *RowsRequestSuite) TestGetInvalidCondition() {
	request := exasol_rest_api.RowsRequest{
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      100,
			},
			ComparisonPredicate: "foo",
		},
	}
	condition, err := request.GetCondition()
	suite.Equal("\"MY_COLUMN\" = 100", condition)
	suite.NoError(err)
}

func (suite *RowsRequestSuite) TestGetInvalidCondition2() {
	request := exasol_rest_api.RowsRequest{
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      nil,
			},
			ComparisonPredicate: "!=",
		},
	}
	condition, err := request.GetCondition()
	suite.Empty(condition)
	suite.EqualError(err, "E-ERA-16: invalid exasol literal type <nil> for value <nil> in the request")
}

func (suite *RowsRequestSuite) TestGetConditionWithDefaultValue() {
	request := exasol_rest_api.RowsRequest{
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	condition, err := request.GetCondition()
	suite.Equal("\"MY_COLUMN\" = 'value'", condition)
	suite.NoError(err)
}

func (suite *RowsRequestSuite) TestValidateSuccess() {
	request := exasol_rest_api.RowsRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.NoError(request.Validate())
}

func (suite *RowsRequestSuite) TestValidateWithoutSchemaName() {
	request := exasol_rest_api.RowsRequest{
		TableName: "MY_TABLE",
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-19: request has some missing parameters. "+
			"Please specify schema name, table name and condition: column name, value")
}

func (suite *RowsRequestSuite) TestValidateWithoutTableName() {
	request := exasol_rest_api.RowsRequest{
		SchemaName: "MY_SCHEMA",
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
				Value:      "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-19: request has some missing parameters. "+
			"Please specify schema name, table name and condition: column name, value")
}

func (suite *RowsRequestSuite) TestValidateWithoutColumnName() {
	request := exasol_rest_api.RowsRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				Value: "value",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-19: request has some missing parameters. "+
			"Please specify schema name, table name and condition: column name, value")
}

func (suite *RowsRequestSuite) TestValidateWithoutColumnValue() {
	request := exasol_rest_api.RowsRequest{
		SchemaName: "MY_SCHEMA",
		TableName:  "MY_TABLE",
		WhereCondition: exasol_rest_api.Condition{
			CellValue: exasol_rest_api.Value{
				ColumnName: "MY_COLUMN",
			},
		},
	}
	suite.EqualError(request.Validate(),
		"E-ERA-19: request has some missing parameters. "+
			"Please specify schema name, table name and condition: column name, value")
}
