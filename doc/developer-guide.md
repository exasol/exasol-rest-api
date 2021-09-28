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