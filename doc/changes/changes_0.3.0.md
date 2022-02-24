# Exasol REST API 0.3.0, released 2022-02-25

Code name: Security- and bug fixes

## Summary

In this release we fixed some security issues in dependencies and added [instructions](../user_guide/user-guide.md#accessing-the-service) how to try out the REST API on the command line.

## Bugfixes

* #53: Fixed build under macOS
* #54: Fixed security issues by upgrading dependencies:
  * [GHSA-qq97-vm5h-rrhg](https://github.com/advisories/GHSA-qq97-vm5h-rrhg)
  * [CVE-2021-43784](https://github.com/advisories/GHSA-v95c-p5hm-xq8f)
  * [GHSA-77vh-xpmg-72qh](https://github.com/advisories/GHSA-77vh-xpmg-72qh)
  * [CVE-2021-30465](https://github.com/advisories/GHSA-c3xm-pvg7-gh7r)

## Upgrades

### Direct Dependencies

* Upgraded `github.com/exasol/error-reporting-go` to `v0.1.1`
* Upgraded `github.com/exasol/exasol-driver-go` to `v0.3.0`
* Upgraded `github.com/gorilla/websocket v1.4.2` to `v1.5.0+incompatible`
* Upgraded `github.com/swaggo/gin-swagger v1.3.2` to `v1.4.1`
* Upgraded `github.com/swaggo/swag v1.7.3` to `v1.8.0`
* Upgraded `github.com/testcontainers/testcontainers-go v0.11.1` to `v.12.0`
* Upgraded `github.com/tidwall/sjson v1.2.3` to `v.1.2.4`
* Upgraded `github.com/ulule/limiter/v3 v3.8.0` to `v3.9.0`

### Indirect Dependencies

* Added `github.com/gofrs/uuid v4.2.0+incompatible`
* Added `github.com/shopspring/decimal v1.3.1`
* Added `go.opentelemetry.io/otel v1.4.1`
* Added `go.opentelemetry.io/otel v1.4.1`
* Added `github.com/go-openapi/jsonreference v0.19.6`
* Added `golang.org/x/text v0.3.7`
* Added `gopkg.in/go-playground/assert.v1 v1.2.1`
* Added `gopkg.in/go-playground/validator.v8 v8.18.2`
* Upgraded `github.com/go-openapi/swag v0.19.15` to `v0.21.1`
* Upgraded `golang.org/x/net v0.0.0-20211013171255-e13a2654a71e` to `v0.0.0-20220127200216-cd36cc0744dd`
* Upgraded `golang.org/x/sys v0.0.0-20211013075003-97ac67df715c` to `v0.0.0-20220209214540-3681064d5158`
* Upgraded `golang.org/x/tools v0.1.7` to `v0.1.9`
* Removed `github.com/Microsoft/go-winio v0.5.2`
* Removed `github.com/containerd/containerd v1.6.0`
* Removed `github.com/docker/distribution v2.8.0+incompatible`
* Removed `github.com/docker/docker v20.10.12+incompatible`
* Removed `github.com/google/go-cmp v0.5.7`
* Removed `github.com/magiconair/properties v1.8.6`
* Removed `github.com/mattn/go-isatty v0.0.14`
* Removed `github.com/moby/sys/mount v0.3.1`
* Removed `github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6`
* Removed `github.com/opencontainers/image-spec v1.0.2`
* Removed `github.com/tidwall/gjson v1.14.0`
* Removed `github.com/ugorji/go v1.2.6`
* Removed `golang.org/x/crypto v0.0.0-20220214200702-86341886e292`
* Removed `google.golang.org/genproto v0.0.0-20220222213610-43724f9ea8cf`
* Removed `github.com/go-playground/validator/v10 v10.10.0`
