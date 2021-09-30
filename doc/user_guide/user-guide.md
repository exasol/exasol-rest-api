# Exasol REST API

## Installation

### Database Side Preparation

* Create a user in the Exasol database for the API service. This user is a service account for Exasol REST API. The service account will be used for all connections, therefore grant **minimum** privileges to it. We strongly recommend to only grant SELECT on required schemas/tables. Please keep in mind that all API users use the same service account.

```sql
CREATE USER api_service_account IDENTIFIED BY secret_password;
GRANT CREATE SESSION TO api_service_account;
GRANT SELECT ON SCHEMA my_schema TO api_service_account;
```

### Configuration File

The application needs a configuration file in yml format. The minimum required parameters are the credentials of the
Exasol service account described above:

```yaml
ExasolUser: "api_service_account"
ExasolPassword: "secret_password"
```

You can also provide additional configurations:

| Property                     |  Default          | Description                                          |
| :--------------------------: | :---------------: | :--------------------------------------------------- |
| api-tokens                   |                   | List of allowed API tokens for authorization.        |
| server-address               |  "localhost:8080" | Address for the server to listen for new connection. |
| exasol-user                  |                   | Name of the Exasol service account.                  |
| exasol-password              |                   | Password of the Exasol service account.              |
| exasol-host                  | "localhost"       | Exasol host.                                         |
| exasol-port                  | 8563              | Exasol port.                                         |
| exasol-websocket-api-version | 2                 | Version of Exasol Websocket API.                     |
| encryption                   | false             | Automatic [Exasol connection encryption][1]. You can enable or disable it. |
| use-tls                      | false             | TLS/SSL verification. Disable it if you want to use a self-signed or invalid certificate (server side).  |

Before starting the application, you need to set an environment variable that points to the properties file:

```
APPLICATION_PROPERTIES_PATH=application-properties.yml
```

### Authorization

Add a list of API tokens to the configuration file. Example:

```yaml
api-tokens:
  - "fwe3cqzE9pQblAYbLFRmxtN03uMgJ2"
  - "nKBwcSyMHr1BnYsV8kiaU0cxNY6iQr"
  - "ubwl5sCao6RHE3iCqe72M6zJc1cHHQ"
```

The tokens must have at least 30 alphanumeric characters.
Only users with the tokens you listed can access secured API endpoints.

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