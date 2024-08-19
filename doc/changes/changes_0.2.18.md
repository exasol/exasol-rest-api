# Exasol Rest api 0.2.18, released 2024-08-19

Code name: Fix vulnerability

## Summary

This release changes the supported version of `go` to be `1.22` or higher and fixes the vulnerability "Authz zero length regression" in test dependency `github.com/docker/docker:v26.0.2` by updating dependencies.

## Security Issues

* #96: Fix Security Issue Authz zero length regression

## Dependency Updates

### Test Dependency Updates

* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.5` to `v0.3.9`
* Updated `github.com/testcontainers/testcontainers-go:v0.29.1` to `v0.32.0`
