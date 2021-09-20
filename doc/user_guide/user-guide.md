# Exasol REST API

## Installation

### Database Side Preparation

* Create a user in the Exasol database for the api service. This user will be used for all connections. Grant **
  minimum** privileges to the user. We strongly recommend to only grant SELECT on required schemas/tables. Please keep
  in mind that all API users uses the same service account.

```sql
CREATE USER api_service_account IDENTIFIED BY secret_password;
GRANT CREATE SESSION TO api_service_account;
GRANT SELECT ON SCHEMA my_schema TO api_service_account;
```

### Configuration File

The application needs a configuration file in yml format. The minimum required parameters are the credentials of the
Exasol service account described above:

```yaml
ExasolUser: "someuser"
ExasolPassword: "somepass"
```

You can also provide additional configurations:

| Property                  |  Default          | Description                                     |
| :-----------------------: | :---------------: | :---------------------------------------------- |
| ApplicationServer         |  "localhost:8080" | Server to start the API application on.         |
| ExasolUser                |                   | Name of the Exasol service user.                |
| ExasolPassword            |                   | Password of the Exasol service user             |
| ExasolHost                | "localhost"       | Exasol host.                                    |
| ExasolPort                | 8563              | Exasol port.                                    |
| ExasolWebsocketApiVersion | 2                 | Version of Exasol Websocket API.                |
| Encryption                | false             | Exasol automatic encryption.                    |
| UseTLS                    | false             | TLS/SSL verification. Disable it if you want to use a self-signed or invalid certificate (server side).  |

Before starting the application, you need to set an environment variable that points to the properties file:

```
APPLICATION_PROPERTIES_PATH=application-properties.yml
```

### Start API Service

* Download the latest executable from our [GitHub repository](https://github.com/exasol/exasol-rest-api/releases).

* Start the service:

On Linux:

```shell
./exasol-rest-api
```