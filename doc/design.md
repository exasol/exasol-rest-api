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

#### Execute Query Request Body
`dsn~execute-query-request-body~1`

ERA accepts the following format of the request body:

```
{
     "sqlText": <string>
 }
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

ERA provide the following endpoint to the API users: `/api/v1/tables [get]`

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

See a response format of [Exasol WebSocker API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

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

ERA provide the following endpoint to the API users: `/api/v1/row [post]`

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

```
{
     "schemaName": <string>,
     "tableName": <string>,
     "row": {
        "<column name>" : "<value>",
        "<column name>" : "<value>",
        ...
     }
 }
```

Covers:

* `req~insert-row~1`

Needs: impl, itest

#### Get Tables Response Body
`dsn~insert-row-response-body~1`

See a response format of [Exasol WebSocker API](https://github.com/exasol/websocket-api/blob/master/docs/commands/executeV1.md).

Covers:

* `req~insert-row~1`

Needs: impl, itest

# Cross-cutting Concerns

# Design Decisions

# Quality Scenarios

# Risks