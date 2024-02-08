FROM golang:1.20-alpine
COPY . /exasol-rest-api
WORKDIR /exasol-rest-api
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN sh ./generate-swagger-docs
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/exasol-rest-api -buildvcs=false
CMD [ "build/exasol-rest-api" ]
