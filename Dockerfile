FROM golang:1.18-alpine
COPY . /exasol-rest-api
WORKDIR /exasol-rest-api
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.7
RUN sh ./generate-swagger-docs
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/exasol-rest-api -buildvcs=false
CMD [ "build/exasol-rest-api" ]