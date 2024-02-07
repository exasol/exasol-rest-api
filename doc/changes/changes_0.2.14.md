# Exasol Rest api 0.2.14, released 2024-02-08

Code name: Fix vulnerabilities, update direct dependencies + go version

## Summary
Fix vulnerabilities, update direct dependencies + go version
Fixed vulnerabilies: 
* github.com/opencontainers/runc : CVE-2024-21626
* github.com/containerd/containerd : CVE-2024-21626

## Security

* #88: fix vulnerabilities / update dependencies

## Dependency Updates

### Compile Dependency Updates

* Updated `golang:1.20` to `1.21`
* Updated `github.com/swaggo/swag:v1.16.2` to `v1.16.3`
