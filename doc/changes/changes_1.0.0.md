# Exasol Rest api 0.3.0, released 2025-05-07

Code name: Support Exasol v8

## Summary

This release adds support for Exasol v8 and refactors the configuration of Exasol certificate validation. See the [user guide section](../user_guide/user-guide.md#encrypted-connection-to-the-exasol-database) for details.

**Breaking Changes:**
* The REST API now only supports TLS encrypted connections to the Exasol database. Unencrypted connections to Exasol 7.1 servers are not supported any more.
* The REST API now uses compression for connections to the Exasol database.
* The REST API does not support the following properties any more:
  * `EXASOL_ENCRYPTION`: This option allowed to enable or deactivate encryption when connecting to an Exasol database. Encryption is now enabled by default and cannot be deactivated.
  * `EXASOL_TLS`: This option allowed to enable or deactivate verification of the Exasol database TLS certificate and supported values `1` and `-1`. We replaced this option with property `EXASOL_VALIDATE_SERVER_CERTIFICATE` that supports values `true` (default) and `false`.
* REST API now additionally supports the following configuration properties:
  * `EXASOL_VALIDATE_SERVER_CERTIFICATE`: Enable (`true`, default) or disable (`false`) verification of the Exasol TLS certificate.
  * `EXASOL_CERTIFICATE_FINGERPRINT`: Expected fingerprint of the Exasol TLS certificate. This is useful when Exasol uses a self-signed certificate.
* Column type metadata for queries does not contain the following fields any more:
  * `characterSet` (values: `"UTF8"`, `"ASCII"`)
  * `withLocalTimeZone` (values: `true`, `false`)
  * `fraction`
  * `srid`
* Configuration property `EXASOL_WEBSOCKET_API_VERSION` is not supported any more and is ignored.
* The REST API now returns the correct error status code 500 (Internal server error) instead of 400 (Bad request) when connection to the Exasol database fails.
* Error messages e.g. for connection failures or query errors are different in the new version.

**Note:**
* For backwards compatibility the REST API returns status code 200 instead of 400 when a query or statement fails e.g. due to syntax error.

## Features

* #87: Add support for Exasol v8

## Dependency Updates

### Compile Dependency Updates

* Updated `golang:1.21` to `1.24.0`
* Updated `github.com/swaggo/swag:v1.16.3` to `v1.16.4`
* Added `github.com/exasol/exasol-driver-go:v1.0.13`
* Updated `github.com/stretchr/testify:v1.9.0` to `v1.10.0`

### Test Dependency Updates

* Updated `github.com/testcontainers/testcontainers-go:v0.32.0` to `v0.37.0`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.9` to `v0.3.11`

### Other Dependency Updates

* Added `toolchain:go1.24.2`
* Removed `github.com/gorilla/websocket:v1.5.3`
