# Exasol REST API

## Installation

### Database Side Preparation

* Create a user in the Exasol database for the API service. This user is a service account for Exasol REST API. The service account will be used for all connections, therefore grant **minimum** privileges to it. We strongly recommend to only grant SELECT on required schemas/tables. Please keep in mind that all API users use the same service account.

```sql
CREATE USER api_service_account IDENTIFIED BY secret_password;
GRANT CREATE SESSION TO api_service_account;
GRANT SELECT ON SCHEMA my_schema TO api_service_account;
```

### Properties

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

Before starting the application, you need to set an environment variable that points to the properties file:

```
APPLICATION_PROPERTIES_PATH=application-properties.yml
```

#### Via Environment Variables

You can set the properties via environment variables. Use the properties' names from the table below.
For the API tokens' value use the following format: `token1,token2,token3,...`

#### Properties Reading Chain

1. The properties from the YAML file are read.
2. The configuration from environment variables are read, they override the configurations from the YAML file.
3. The default values are added if they were not added in the previous steps.

#### All Available Properties

| Property                     |  Default        | Description                                          |
| :--------------------------: | :-------------: | :--------------------------------------------------- |
| API_TOKENS                   |                 | List of allowed API tokens for authorization.        |
| SERVER_ADDRESS               |  "0.0.0.0:8080" | Address for the server to listen for new connection. |
| EXASOL_USER                  |                 | Name of the Exasol service account.                  |
| EXASOL_PASSWORD              |                 | Password of the Exasol service account.              |
| EXASOL_HOST                  | "localhost"     | Exasol host.                                         |
| EXASOL_PORT                  | 8563            | Exasol port.                                         |
| EXASOL_WEBSOCKET_API_VERSION | 2               | Version of Exasol Websocket API.                     |
| EXASOL_ENCRYPTION            | false           | Automatic [Exasol connection encryption][1]. You can enable or disable it. |
| EXASOL_TLS                   | false           | TLS/SSL verification. Disable it if you want to use a self-signed or invalid certificate (server side).  |

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

Windows Hint: If you start the application via a double-click on the file, when the application crashes you it exits immediately. It means you won't see any error messages. So we recommend starting the application via cmd.

### Accessing the Service

You can access the service on the host and port you specified. For the default values: `http://localhost:8080/api/v1/<endpoint here>`.

You can also access the Swagger documentation. Here is an example with the default values: `http://localhost:8080/swagger/index.html`

[1]: https://community.exasol.com/t5/database-features/database-connection-encryption-at-exasol/ta-p/2259

### Rate Limitation

The service allows 30 requests per minute for all API endpoints. The limitation is based on the sender's IP address.