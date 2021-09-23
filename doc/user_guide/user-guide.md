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
APPLICATION_PROPERTIES_PATH=application_properties.yml
```

### Start API Service

* Download the latest executable from our [GitHub repository](https://github.com/exasol/exasol-rest-api/releases).

* Start the service:

On Linux:

```shell
./exasol-rest-api
```

[1]: https://community.exasol.com/t5/database-features/database-connection-encryption-at-exasol/ta-p/2259