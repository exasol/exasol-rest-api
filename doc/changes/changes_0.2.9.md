# Exasol Rest api 0.2.9, released 2023-05-08

Code name: Fix vulnerabilities in test dependency

## Summary

This release upgrades test dependency `github.com/docker/docker` to fix vulnerabilities CVE-2023-28840, CVE-2023-28841 and CVE-2023-28842.

## Features

* #75: Fixed vulnerabilities in test dependency

## Dependency Updates

### Compile Dependency Updates

* Updated `github.com/swaggo/swag:v1.8.12` to `v1.16.1`
* Updated `github.com/exasol/error-reporting-go:v0.1.1` to `v0.2.0`

### Test Dependency Updates

* Updated `github.com/testcontainers/testcontainers-go:v0.19.0` to `v0.20.0`
* Updated `github.com/exasol/exasol-driver-go:v0.4.7` to `v1.0.0`
