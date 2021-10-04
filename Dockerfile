FROM golang:1.16-alpine
COPY . /exasol-rest-api
WORKDIR /exasol-rest-api
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN export PATH=$(go env GOPATH)/bin:$PATH
RUN sh generate-swagger-docs
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/exasol-rest-api
ENV APPLICATION_PROPERTIES_PATH=/exasol-rest-api/application-properties.yml
CMD [ "build/exasol-rest-api" ]