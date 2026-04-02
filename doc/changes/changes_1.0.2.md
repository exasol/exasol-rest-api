# Exasol Rest api 1.0.2, released 2026-04-02

Code name: Fix vulnerabilities in dependencies

## Summary

This release fixes the following vulnerabilities in dependencies.

### Govuln Findings
* [GO-2025-4135](https://pkg.go.dev/vuln/GO-2025-4135): Malformed constraint may cause denial of service in `golang.org/x/crypto/ssh/agent`
* [GO-2025-4233](https://pkg.go.dev/vuln/GO-2025-4233): HTTP/3 QPACK Header Expansion DoS in `github.com/quic-go/quic-go`
* [GO-2025-4134](https://pkg.go.dev/vuln/GO-2025-4134): Unbounded memory consumption in `golang.org/x/crypto/ssh`

### Dependabot Alerts
* [#50](https://github.com/exasol/exasol-rest-api/security/dependabot/50): Moby has AuthZ plugin bypass when provided oversized request bodies
* [#49](https://github.com/exasol/exasol-rest-api/security/dependabot/49): Moby has an Off-by-one error in its plugin privilege validation
* [#45](https://github.com/exasol/exasol-rest-api/security/dependabot/45): golang.org/x/crypto/ssh allows an attacker to cause unbounded memory consumption
* [#48](https://github.com/exasol/exasol-rest-api/security/dependabot/48): quic-go HTTP/3 QPACK Header Expansion DoS
* [#46](https://github.com/exasol/exasol-rest-api/security/dependabot/46): golang.org/x/crypto/ssh/agent vulnerable to panic if message is malformed due to out of bounds read

## Security

* #106: Fix vulnerabilities in dependencies
