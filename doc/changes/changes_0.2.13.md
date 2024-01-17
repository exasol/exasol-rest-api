# Exasol Rest API 0.2.13, released 2024-01-17

Code name: Fix CVE-2023-48795 in `golang.org/x/crypto`

## Summary

This release fixes CVE-2023-48795 in transitive dependency `golang.org/x/crypto`

## Security

* #85: Fixed CVE-2023-48795 in `golang.org/x/crypto`

## Dependency Updates

### Compile Dependency Updates

* Updated `github.com/gorilla/websocket:v1.5.0` to `v1.5.1`

### Test Dependency Updates

* Updated `github.com/testcontainers/testcontainers-go:v0.25.0` to `v0.27.0`
* Updated `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.4` to `v0.3.5`
