# Developer Guide

## Swagger

### API Version

Go project doesn't have a defined version. We provide version in the `main.go` file to display on the Swagger page. 
Please don't forget to change that version.

### Generate Swagger Documentation Locally

* Go to the root directory of the project.
* Install `swag`:

```shell
 go get -v -u github.com/swaggo/swag/cmd/swag
```

* Run the following script to generate Swagger documentation:

```shell
bash generate-swagger-docs 
 ```

Hint: If you installed swag, but the script doesn't work because swag is missing, try setting the following environment variable:

```shell
export PATH=$(go env GOPATH)/bin:$PATH
```

### Building Docker Image Locally

We are planning to deliver a docker image via a registry.
If you need to develop or test the image, you can create it locally.

To start a docker container:

1. Clone the REST API repository.
2. Switch to the repository root directory: `cd exasol-rest-api`
3. Build the docker image: `docker build --tag my-rest-api-image:latest .`
4. Start a container: `docker run --name test-api-container --env EXASOL_USER=sys --env EXASOL_PASSWORD=secret --env API_TOKENS=token1,toekn2 -p 8080:8080 my-rest-api-image:latest`