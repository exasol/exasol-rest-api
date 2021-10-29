# Introduction

## Acknowledgments

This document's section structure is derived from the "[arc42](https://arc42.org/)" architectural template by Dr. Gernot Starke, Dr. Peter Hruschka.

## Terms and Abbreviations

<dl>
    <dt>ERA</dt><dd>Exasol REST API</dd>
    <dt>Exasol admin</dt><dd>A person with an Exasol admin account</dd>
</dl>

## Requirement Overview

Please refer to the [System Requirement Specification](system-requirements.md) for user-level requirements.

# Building Blocks

This section introduces the building blocks of the software. Together those building blocks make up the big picture of the software structure.

## Proxy Service

An application that works as a proxy between an API user and [Exasol WebSockets API](https://github.com/exasol/websocket-api). 

# Runtime

This section describes the runtime behavior of the software.

## Authentication

## Proxy Service

### Communicate with Database
`dsn~communicate-with-database~1`

ERA uses [Exasol WebSockets API](https://github.com/exasol/websocket-api) to proxy the communication with the Exasol database.

Covers:

* `req~execute-query~1`

Needs: impl, itest

### Service Account
`dsn~service-account~1`

Exasol admins create a service account for the ERA proxy service.

Covers:

* `req~execute-query~1`

Needs: impl, itest

### Service Credentials
`dsn~service-credentials~1`

Exasol admins configure the proxy service with the service account credentials.

Covers:

* `req~execute-query~1`

Needs: impl, itest

## REST API Endpoints

### Execute Query

#### Execute Query Endpoint 
`dsn~execute-query-endpoint~1`

ERA provide the following endpoint to the API users: `/api/v1/query [get]`

Covers:

* `req~execute-query~1`

Needs: impl, itest

#### Execute Query Headers
`dsn~execute-query-headers~1`

The endpoint requires `Authorization` header with an API token to handle requests.

Covers:

* `req~execute-query~1`

Needs: impl, itest

#### Execute Query Request Parameters
`dsn~execute-query-request-parameters~1`

The endpoint accepts the following path parameter:

```
<endpoint>/<query string>
```

Covers:

* `req~execute-sql-query~1`

Needs: impl, itest

#### Execute Query Response Body
`dsn~execute-query-response-body~1`

See a response format of [Exasol WebSocker API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

Covers:

* `req~execute-sql-query~1`

Needs: impl, itest

#### Query Results Limitation
`dsn~query-results-limitation~1`

The result set has 1000 rows or fewer.

Covers:

* `req~execute-sql-query~1`

Needs: impl, itest

### Get Tables

#### Get Tables Endpoint
`dsn~get-tables-endpoint~1`

ERA provides the following endpoint to the API users: `/api/v1/tables [get]`

Covers:

* `req~get-tables~1`

Needs: impl, itest

#### Get Tables Headers
`dsn~get-tables-headers~1`

The endpoint requires `Authorization` header with an API token to handle requests.

Covers:

* `req~get-tables~1`

Needs: impl, itest

#### Get Tables Response Body
`dsn~get-tables-response-body~1`

See a response format of [Exasol WebSocket API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

Covers:

* `req~get-tables~1`

Needs: impl, itest

#### Get Tables Results Limitation
`dsn~get-tables-results-limitation~1`

The result set has 1000 rows or fewer.

Covers:

* `req~get-tables~1`

Needs: impl, itest

### Insert Row

#### Insert Row Endpoint
`dsn~insert-row-endpoint~1`

ERA provides the following endpoint to the API users: `/api/v1/row [post]`

Covers:

* `req~insert-row~1`

Needs: impl, itest

#### Insert Row Headers
`dsn~insert-row-headers~1`

The endpoint requires `Authorization` header with an API token to handle requests.

Covers:

* `req~insert-row~1`

Needs: impl, itest

#### Insert Row Request Body
`dsn~insert-row-request-body~1`

ERA accepts the following format of the request body:

```json
{
     "schemaName": <string>,
     "tableName": <string>,
     "row": [
     	{
        	"columnName": <string>,
     		"value": <value>
     	},
     	{
        	"columnName": <string>,
     		"value": <value>
     	},
     	...
     ]
 }
```

Covers:

* `req~insert-row~1`

Needs: impl, itest

#### Insert Row Response Body
`dsn~insert-row-response-body~1`

See a response format of [Exasol WebSocket API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

Covers:

* `req~insert-row~1`

Needs: impl, itest

### Delete Rows

#### Delete Rows Endpoint
`dsn~delete-rows-endpoint~1`

ERA provide the following endpoint to the API users: `/api/v1/rows [delete]`

Covers:

* `req~delete-rows~1`

Needs: impl, itest

#### Delete Rows Headers
`dsn~delete-rows-headers~1`

The endpoint requires `Authorization` header with an API token to handle requests.

Covers:

* `req~delete-rows~1`

Needs: impl, itest

#### Delete Rows Request Body
`dsn~delete-rows-request-body~1`

ERA accepts the following format of the request body:

```json
{
     "schemaName": <string>,
     "tableName": <string>,
     "condition": {
        "value": {
            "columnName": <string>,
            "value": <value>,
        },
        "comparisonPredicate": "= or != or < or <= or > or >="
     }
 }
```

Covers:

* `req~delete-rows~1`

Needs: impl, itest

#### Delete Rows Response Body
`dsn~delete-rows-response-body~1`

See a response format of [Exasol WebSocket API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

Covers:

* `req~delete-rows~1`

Needs: impl, itest

### Get Rows

#### Get Rows Endpoint
`dsn~get-rows-endpoint~1`

ERA provides the following endpoint to the API users: `/api/v1/rows [get]`

Covers:

* `req~get-rows~1`

Needs: impl, itest

#### Get Rows Headers
`dsn~get-rows-headers~1`

The endpoint requires `Authorization` header with an API token to handle requests.

Covers:

* `req~get-rows~1`

Needs: impl, itest

#### Get Rows Request Parameters
`dsn~get-rows-request-parameters~1`

The endpoint accepts the following query string parameters:

```
<endpoint>?schemaName=<schema name>&tableName=<table name>&columnName=<column name>&value=<value>&valueType=<string/bool/int/float>&comparisonPredicate=<comparison predicate>
```

Covers:

* `req~get-rows~1`

Needs: impl, itest

#### Get Rows Response Body
`dsn~get-rows-response-body~1`

See a response format of [Exasol WebSocket API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

Covers:

* `req~get-rows~1`

Needs: impl, itest

### Update Rows

#### Update Rows Endpoint
`dsn~update-rows-endpoint~1`

ERA provides the following endpoint to the API users: `/api/v1/rows [put]`

Covers:

* `req~update-rows~1`

Needs: impl, itest

#### Update Rows Headers
`dsn~update-rows-headers~1`

The endpoint requires `Authorization` header with an API token to handle requests.

Covers:

* `req~update-rows~1`

Needs: impl, itest

#### Update Rows Request Body
`dsn~update-rows-request-body~1`

ERA accepts the following format of the request body:

```json
{
     "schemaName": <string>,
     "tableName": <string>,
     "row": [
     	{
        	"columnName": <string>,
     		"value": <value>
     	},
     	{
        	"columnName": <string>,
     		"value": <value>
     	},
     	...
     ],
     "condition": {
        "value": {
            "columnName": <string>,
            "value": <value>,
        },
        "comparisonPredicate": "= or != or < or <= or > or >="
     }
 }
```

Covers:

* `req~update-rows~1`

Needs: impl, itest

#### Update Rows Response Body
`dsn~update-rows-response-body~1`

See a response format of [Exasol WebSocket API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

Covers:

* `req~update-rows~1`

Needs: impl, itest

# Cross-cutting Concerns

# Design Decisions

# Quality Scenarios

# Risks