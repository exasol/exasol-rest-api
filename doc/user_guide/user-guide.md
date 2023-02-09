# Exasol REST API

Exasol REST API is an extension for the Exasol database that provides the ability to interact with the database via REST API endpoints. The REST API is required by the [Power Apps Connector for Exasol](https://github.com/exasol/power-apps-connector).

## Installation

### Deployment

The REST API needs to be installed on a internet-facing server machine/VM that also has connectivity to the database environment.

The REST API is stateless, so you can install it on multiple nodes and put a HTTP load balancer in front.

We recommend to use HTTPS/TLS, see the [available options below](#using-secured-connection-https).

### Database Side Preparation

* Create a user in the Exasol database for the API service. This user is a service account for Exasol REST API. The service account will be used for all connections, therefore grant **minimum** privileges to it. We strongly recommend to only grant SELECT on required schemas/tables. Please keep in mind that all API users use the same service account.

```sql
CREATE USER api_service_account IDENTIFIED BY secret_password;
GRANT CREATE SESSION TO api_service_account;
GRANT SELECT ON SCHEMA my_schema TO api_service_account;
```

### Configuration

#### Via YAML Configuration File

The application needs a configuration file in yml format. The minimum required parameters are the credentials of the
Exasol service account described above and a list of API tokens (see [Authorization](#authorization)):

```yaml
EXASOL_USER: "api_service_account"
EXASOL_PASSWORD: "secret_password"
API_TOKENS:
  - "abc"
  - "bca"
  - "cab"
```

Please be aware of the API token length: only 30 or more characters are allowed.

Before starting the application, you need to 
- set an environment variable that points to the properties file:

```
APPLICATION_PROPERTIES_PATH=application-properties.yml
```
OR

- define the path to this properties path in a CLI flag called application-properties-path:
On Linux:

```shell
./exasol-rest-api -application-properties-path='<PATH>'
```

On Windows: open a command prompt and start the service from the prompt:  

```cmd
path\to\file\exasol-rest-api-x86-64.exe -application-properties-path='<PATH>'
```


#### Via Environment Variables

You can set the properties via environment variables. Use the properties' names from the table below.
For the API tokens' value use the following format: `token1,token2,token3,...`

#### Properties Reading Chain

1. The properties from the YAML file are read.
2. The configuration from environment variables are read, they override the configurations from the YAML file.
3. The default values are added if they were not added in the previous steps.

#### All Available Properties

| Property                     |    Default     | Description                                                  |
| :--------------------------- | :------------: | :----------------------------------------------------------- |
| API_TOKENS                   |                | List of allowed API tokens for authorization.                |
| SERVER_ADDRESS               | "0.0.0.0:8080" | Address for the server to listen for new connection.         |
| EXASOL_USER                  |                | Name of the Exasol service account.                          |
| EXASOL_PASSWORD              |                | Password of the Exasol service account.                      |
| EXASOL_HOST                  |  "localhost"   | Exasol host.                                                 |
| EXASOL_PORT                  |      8563      | Exasol port.                                                 |
| EXASOL_WEBSOCKET_API_VERSION |       2        | Version of Exasol Websocket API.                             |
| EXASOL_ENCRYPTION            |       1        | Automatic [Exasol connection encryption][1]. Use 1 to enable it and -1 to disable. |
| EXASOL_TLS                   |       1        | Database TLS/SSL verification. Disable it if you want to use a self-signed or invalid certificate (server side). Use 1 to enable it and -1 to disable. |
| API_TLS                      |     false      | Enable API TLS/SSL.                                          |
| API_TLS_PKPATH               |                | Path of the private key file.                                |
| API_TLS_CERTPATH             |                | Path of the certificate file.                                |

[1]: https://community.exasol.com/t5/database-features/database-connection-encryption-at-exasol/ta-p/2259

### Authorization

Add a list of API tokens to the configuration file (find an example above).
The tokens must have **at least 30 alphanumeric** characters. Only users with the tokens you listed can access secured API endpoints.

### Start API Service

* Download the latest executable from our [GitHub repository](https://github.com/exasol/exasol-rest-api/releases).

* Start the service:

On Linux:

```shell
./exasol-rest-api
```

On Windows: open a command prompt and start the service from the prompt:

```shell
path\to\file\exasol-rest-api-x86-64.exe
```

Windows Hint: If you start the application via a double-click on the file, when the application crashes it exits immediately. That means you won't see any error messages. So we recommend starting the application via CLI.

### Accessing the Service

You can access the service on the host and port you specified. For the default values: `http://localhost:8080/api/v1/<endpoint here>`.

The Swagger documentation of the API is available at `http://localhost:8080/swagger/index.html`. After authenticating with a token you can try out the requests.

You can also send requests via command line:

```sh
token="<token>"
base_url="http://localhost:8080/api/v1"
# List all tables
curl --header "Authorization: $token" "$base_url/tables"
# Execute a query. Spaces and other characters must be url-encoded.
curl --header "Authorization: $token" "$base_url/query/select%201"
```

See the [design document](../design.md) for payload examples.

### Rate Limitation

The service allows 30 requests per minute for all API endpoints. The limitation is based on the sender's IP address.

### Using Secured Connection (HTTPS)

We strongly recommend using TLS/HTTPS when deploying the API service.

#### Enable TLS Within the API Itself

You can enable HTTPS within the service itself using the `API_TLS` (set it to true to enable TLS), `API_TLS_PKPATH` (path to `private key.pem` file) and `API_TLS_CERTPATH` (path to `certificate.pem` file) configuration properties.

These entries would look like this in your configuration `.yml` file:

> API_TLS: true
>
> API_TLS_PKPATH: "C:\\\tls\\\private.key.pem"
>
> API_TLS_CERTPATH: "C:\\\tls\\\domain.cert.pem"

You might also want to change the listening port to port 443 (the default for SSL for browsers and many other applications). You can do this by changing the `SERVER_ADDRESS` to  "0.0.0.0:443".

> SERVER_ADDRESS : "0.0.0.0:443"

#### Using a Proxy and TLS Termination

Another way is setting up a publicly accessible proxy server that [terminates SSL](https://en.wikipedia.org/wiki/TLS_termination_proxy) and forwards the request to the API service.

Here are a few examples of services you could use:

* For a local setup: [NGINX](https://docs.nginx.com/nginx/admin-guide/security-controls/terminating-ssl-http/)
* For AWS setup: [Application Load Balancer](https://aws.amazon.com/elasticloadbalancing/)
* For Azure setup: [Azure Front Door](https://docs.microsoft.com/en-us/azure/frontdoor/front-door-overview)
