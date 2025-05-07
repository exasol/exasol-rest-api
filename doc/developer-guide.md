# Developer Guide

## Swagger

### API Version

Go project doesn't have a defined version. We provide version in the `main.go` file to display on the Swagger page.
Please don't forget to change that version.

### Generate Swagger Documentation Locally

1. Go to the root directory of the project.
2. Install `swag`:

    ```shell
    go install github.com/swaggo/swag/cmd/swag@v1.16.4
    ```

3. Run the following script to generate Swagger documentation:

    ```shell
    ./generate-swagger-docs
    ```

### Building Docker Image Locally

We are planning to deliver a docker image via a registry. If you need to develop or test the image, you can create it locally.

To start a docker container:

1. Clone the REST API repository.
2. Switch to the repository root directory: `cd exasol-rest-api`
3. Build the docker image:

    ```shell
    docker buildx build --tag rest-api-test-image:latest .
    ```

4. Generate one or more random tokens with at least 30 characters, e.g. by calling `uuidgen`:

    ```shell
    token1=$(uuidgen --random)
    token2=$(uuidgen --random)
    echo token1=$token1
    echo token2=$token2
    ```

5. Start a container:

    ```shell
    docker run --name test-api-container \
               --env EXASOL_USER=sys \
               --env EXASOL_PASSWORD=secret \
               --env API_TOKENS=$token1,$token2 \
               -p 8080:8080 \
               my-rest-api-image:latest
    ```
