# Exasol REST API 0.2.3, released 2022-03-30

Code name: Security fixes

## Summary

In this release we fixed some security issues in dependencies.

## Bugfixes

* #56: Fixed security issue by upgrading dependencies: [CVE-2022-23648](https://github.com/advisories/GHSA-crp2-qrr5-8pq7)

## Upgrades

### Direct Dependencies

* Upgraded `golang 1.16` to `1.18`
* Upgraded `github.com/swaggo/swag v1.8.0` to `v1.8.1`
* Upgraded `github.com/exasol/exasol-driver-go v0.3.0` to `v0.3.1`
* Upgraded `github.com/stretchr/testify v1.7.0` to `v1.7.1`
* Upgraded `github.com/ulule/limiter/v3 v3.9.0` to `v3.10.0`
* Upgraded `github.com/testcontainers/testcontainers-go v0.12.0` to `v0.12.1-0.20220216090119-c0c2f90f591a`

### Indirect Dependencies

* Upgraded `github.com/Azure/go-ansiterm` to `v0.0.0-20210617225240-d185dfc1b5a1`
* Upgraded `github.com/KyleBanks/depth` to `v1.2.1`
* Upgraded `github.com/Microsoft/go-winio` to `v0.5.2`
* Upgraded `github.com/Microsoft/hcsshim` to `v0.9.2`
* Upgraded `github.com/PuerkitoBio/purell` to `v1.1.1`
* Upgraded `github.com/PuerkitoBio/urlesc` to `v0.0.0-20170810143723-de5bf2ad4578`
* Upgraded `github.com/cenkalti/backoff` to `v2.2.1+incompatible`
* Upgraded `github.com/cenkalti/backoff/v4` to `v4.1.2`
* Upgraded `github.com/containerd/cgroups` to `v1.0.3`
* Upgraded `github.com/containerd/containerd` to `v1.6.2`
* Upgraded `github.com/davecgh/go-spew` to `v1.1.1`
* Upgraded `github.com/docker/distribution` to `v2.8.1+incompatible`
* Upgraded `github.com/docker/docker` to `v20.10.14+incompatible`
* Upgraded `github.com/docker/go-connections` to `v0.4.0`
* Upgraded `github.com/docker/go-units` to `v0.4.0`
* Upgraded `github.com/gin-contrib/sse` to `v0.1.0`
* Upgraded `github.com/go-openapi/jsonpointer` to `v0.19.5`
* Upgraded `github.com/go-openapi/jsonreference` to `v0.19.6`
* Upgraded `github.com/go-openapi/spec` to `v0.20.4`
* Upgraded `github.com/go-openapi/swag` to `v0.21.1`
* Upgraded `github.com/go-playground/locales` to `v0.14.0`
* Upgraded `github.com/go-playground/universal-translator` to `v0.18.0`
* Upgraded `github.com/go-playground/validator/v10` to `v10.10.1`
* Upgraded `github.com/gogo/protobuf` to `v1.3.2`
* Upgraded `github.com/golang/groupcache` to `v0.0.0-20210331224755-41bb18bfe9da`
* Upgraded `github.com/golang/protobuf` to `v1.5.2`
* Upgraded `github.com/google/go-cmp` to `v0.5.7`
* Upgraded `github.com/google/uuid` to `v1.3.0`
* Upgraded `github.com/josharian/intern` to `v1.0.0`
* Upgraded `github.com/json-iterator/go` to `v1.1.12`
* Upgraded `github.com/leodido/go-urn` to `v1.2.1`
* Upgraded `github.com/magiconair/properties` to `v1.8.6`
* Upgraded `github.com/mailru/easyjson` to `v0.7.7`
* Upgraded `github.com/mattn/go-isatty` to `v0.0.14`
* Upgraded `github.com/moby/sys/mount` to `v0.3.1`
* Upgraded `github.com/moby/sys/mountinfo` to `v0.6.0`
* Upgraded `github.com/moby/term` to `v0.0.0-20210619224110-3f7ff695adc6`
* Upgraded `github.com/modern-go/concurrent` to `v0.0.0-20180306012644-bacd9c7ef1dd`
* Upgraded `github.com/modern-go/reflect2` to `v1.0.2`
* Upgraded `github.com/morikuni/aec` to `v1.0.0`
* Upgraded `github.com/opencontainers/go-digest` to `v1.0.0`
* Upgraded `github.com/opencontainers/image-spec` to `v1.0.2`
* Upgraded `github.com/opencontainers/runc` to `v1.1.1`
* Upgraded `github.com/pkg/errors` to `v0.9.1`
* Upgraded `github.com/pmezard/go-difflib` to `v1.0.0`
* Upgraded `github.com/sirupsen/logrus` to `v1.8.1`
* Upgraded `github.com/tidwall/gjson` to `v1.14.0`
* Upgraded `github.com/tidwall/match` to `v1.1.1`
* Upgraded `github.com/tidwall/pretty` to `v1.2.0`
* Upgraded `github.com/ugorji/go/codec` to `v1.2.7`
* Upgraded `go.opencensus.io` to `v0.23.0`
* Upgraded `golang.org/x/crypto` to `v0.0.0-20220321153916-2c7772ba3064`
* Upgraded `golang.org/x/net` to `v0.0.0-20220325170049-de3da57026de`
* Upgraded `golang.org/x/sys` to `v0.0.0-20220330033206-e17cdc41300f`
* Upgraded `golang.org/x/text` to `v0.3.7`
* Upgraded `golang.org/x/tools` to `v0.1.10`
* Upgraded `google.golang.org/genproto` to `v0.0.0-20220329172620-7be39ac1afc7`
* Upgraded `google.golang.org/grpc` to `v1.45.0`
* Upgraded `google.golang.org/protobuf` to `v1.28.0`
* Upgraded `gopkg.in/yaml.v2` to `v2.4.0`
