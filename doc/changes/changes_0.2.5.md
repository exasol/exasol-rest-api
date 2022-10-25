# exasol-rest-api 0.2.5, released 2022-10-25

Code name: Fix vulnerabilities in dependencies

## Summary

In this release we updated the dependency `pkg:golang/golang.org/x/text` from v0.3.7 to v0.4.0 in order to fix CVE-2022-32149.

## Features

* #64: Fixed vulnerabilities in dependencies
## Dependency Updates

### Compile Dependency Updates

* Updated `github.com/tidwall/sjson:v1.2.4` to `v1.2.5`
* Updated `github.com/swaggo/swag:v1.8.4` to `v1.8.7`
* Updated `github.com/swaggo/gin-swagger:v1.5.1` to `v1.5.3`

### Test Dependency Updates

* Updated `github.com/stretchr/testify:v1.8.0` to `v1.8.1`
* Updated `github.com/testcontainers/testcontainers-go:v0.13.0` to `v0.15.0`
* Updated `github.com/exasol/exasol-driver-go:v0.3.1` to `v0.4.6`
