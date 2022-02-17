FROM golang:1.17-alpine
COPY . /exasol-rest-api
WORKDIR /exasol-rest-api
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.9-p1
RUN sh ./generate-swagger-docs
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/exasol-rest-api
CMD [ "build/exasol-rest-api" ]