# Exasol REST API 1.0.3, released 2026-09-??

Code name: Fix float row filters

## Summary

This release fixes float filters for the `GET /api/v1/rows` endpoint.

**Note:** Starting with this release, the Rest API is not tested with Exasol version 7.1 any more. Only the latest LTS release 2025.1.x and 8.29.x and the latest release 2026.1.1 are supported. Exasol 7.1 is officially past end-of-life since 2026-06-30 (see https://docs.exasol.com/db/latest/planning/life_cycle/life_cycle_policy.htm).

## Bugfixes

* #111: Fix float row filters

## Dependency Updates

### Compile Dependency Updates

* Updated `golang:1.25.0` to `1.26.0`
* Updated `github.com/exasol/error-reporting-go:v0.2.0` to `v0.2.1`
* Updated `github.com/exasol/exasol-driver-go:v1.0.16` to `v1.1.0`
* Updated `github.com/stretchr/testify:v1.11.1` to `v1.12.1`

### Test Dependency Updates

* Updated `github.com/testcontainers/testcontainers-go:v0.41.0` to `v0.44.0`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v1.0.0` to `v1.0.1`

### Other Dependency Updates

* Updated `toolchain:go1.26.1` to `go1.27.1`
