# Exasol Rest api 0.2.8, released 2023-04-03

Code name: Fix vulnerabilities in test dependencies

## Summary

This release fixes the following vulnerabilities in test dependencies:
* CVE-2023-27561
* CVE-2023-28642
* CVE-2023-25809

## Bugfixes

* #72: Fixed dependabot warning about vulnerabilities in test dependencies

## Dependency Updates

### Compile Dependency Updates

* Updated `github.com/ulule/limiter/v3:v3.10.0` to `v3.11.1`
* Updated `github.com/gin-gonic/gin:v1.8.2` to `v1.9.0`
* Updated `github.com/swaggo/swag:v1.8.7` to `v1.8.12`
* Updated `github.com/swaggo/files:v0.0.0-20220728132757-551d4a08d97a` to `v1.0.1`
* Updated `github.com/swaggo/gin-swagger:v1.5.3` to `v1.6.0`

### Test Dependency Updates

* Updated `github.com/stretchr/testify:v1.8.1` to `v1.8.2`
* Updated `github.com/testcontainers/testcontainers-go:v0.15.0` to `v0.19.0`
* Updated `github.com/exasol/exasol-driver-go:v0.4.6` to `v0.4.7`
