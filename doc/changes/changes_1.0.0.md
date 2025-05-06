# Exasol Rest api 0.3.0, released 2025-05-??

Code name: Support Exasol v8

## Summary

This release adds support for Exasol v8.

**Breaking Changes:**
* Column type metadata for queries does not contain the following fields any more:
  * `characterSet`
  * `withLocalTimeZone`
  * `fraction`
  * `srid`
* Configuration property `EXASOL_WEBSOCKET_API_VERSION` is not supported any more and is ignored.
* Error messages e.g. for connection failures or query errors are different in the new version.

## Features

* #87: Add support for Exasol v8

## Dependency Updates

### Compile Dependency Updates

* Updated `golang:1.21` to `1.24.0`
* Updated `github.com/swaggo/swag:v1.16.3` to `v1.16.4`

### Test Dependency Updates

* Updated `github.com/testcontainers/testcontainers-go:v0.32.0` to `v0.37.0`
* Updated `github.com/stretchr/testify:v1.9.0` to `v1.10.0`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.9` to `v0.3.11`

### Other Dependency Updates

* Added `toolchain:go1.24.1`
