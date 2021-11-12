# System Requirement Specification for Exasol REST API

## Introduction

Exasol REST API (ERA) is an extension for the Exasol database that provides the ability to interact with the database via REST API endpoints.

## About This Document

### Target Audience

The target audience are end-users, requirement engineers, software designers and quality assurance. See section ["Stakeholders"](#stakeholders) for more details.

### Goal

The Exasol REST API's main goal is to enable attaching 3rd party products to Exasol that require a REST-compliant web interface.

### Quality Goals

Exasol REST API's main quality goals are in descending order of importance:

* Standard Compliance
* Usability
* Security
* Performance

## Stakeholders

When reading this section please remember that the listed stakeholders are roles, not people! It is not uncommon in software projects that the same person fulfills multiple roles.

### API Users

People who use the Exasol REST API to interact with the Exasol database.

### Terms and Abbreviations

The following list gives you an overview of terms and abbreviations commonly used in OFT documents.

* ERA - Exasol REST API

## Features

Features are the highest level requirements in this document that describe the main functionality of ERA.

### Exasol REST Endpoints
`feat~exasol-rest-endpoints~1`

ERA provides REST API endpoints that allow API Users to interact with the Exasol database. 

Needs: req

## Functional Requirements

### Communication with Exasol
`req~communication-with-exasol~1`

ERA allows users communicate with Exasol database: send requests and receive responses.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Execute Query
`req~execute-query~1`

API users can execute queries via REST API and receive results back.

Rationale:

Query is a request to access data from a database. If we support queries - we cover the main part of the database read functionality.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Get Tables
`req~get-tables~1`

API users can get a list of table names by a schema name.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Insert Row
`req~insert-row~1`

API users can insert a single row into a table.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Delete Rows
`req~delete-rows~1`

API users can delete rows from a table based on a condition.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Get Rows
`req~get-rows~1`

API users can get rows from a table based on a condition.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Update Rows
`req~update-rows~1`

API users can update rows from a table based on a condition.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Execute Statement
`req~execute-statement~1`

API users can execute statements via REST API.

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn

### Support JSON Request and Response Format
`req~support-json-request-and-response-format~1`

ERA supports a JSON request and response format.

Rationale:

JSON is the most common format for sending and receiving data through a REST API

Covers:

* [feat~exasol-rest-endpoints~1](#exasol-rest-endpoints)

Needs: dsn