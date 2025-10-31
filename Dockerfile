FROM golang:1.25.3-alpine as builder

RUN mkdir /exasol-rest-api
RUN addgroup -S rest-api-user \
 && adduser -S -u 10000 -g rest-api-user rest-api-user
USER rest-api-user

WORKDIR /exasol-rest-api
COPY go.mod go.sum .

RUN go mod download \
 && go install github.com/swaggo/swag/cmd/swag@v1.16.6

COPY . .
USER root
RUN chown -R rest-api-user:rest-api-user /exasol-rest-api
USER rest-api-user
RUN sh ./generate-swagger-docs
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/exasol-rest-api -buildvcs=false

FROM scratch AS final
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
USER rest-api-user

WORKDIR /exasol-rest-api
COPY --from=builder /exasol-rest-api/build/exasol-rest-api /exasol-rest-api/build/exasol-rest-api
CMD [ "build/exasol-rest-api" ]
