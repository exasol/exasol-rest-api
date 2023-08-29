# Exasol Rest api 0.2.11, released 2023-08-29

Code name: Update Dependencies on top of 0.2.10

## Summary

This release fixes vulnerability CVE-2023-3978 in dependency `pkg:golang/golang.org/x/net` by upgrading it to the latest version.

## Security

* #80: Fixed vulnerability CVE-2023-3978 in dependency `pkg:golang/golang.org/x/net`

## Dependency Updates

### Test Dependency Updates

* Updated `github.com/testcontainers/testcontainers-go:v0.20.1` to `v0.23.0`
* Added `github.com/exasol/exasol-test-setup-abstraction-server/go-client:v0.3.3`

### Other Dependency Updates

* Removed `github.com/exasol/exasol-driver-go:v1.0.0`
