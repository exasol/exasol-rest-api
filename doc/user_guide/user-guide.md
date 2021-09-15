# Exasol REST API

## Database Side Preparation

* Create a user in the Exasol database for the api service. This user will be used for all connections. Grant **
  minimum** privileges to the user. We strongly recommend to only grant SELECT on required schemas/tables. Please keep
  in mind that all API users uses the same service account.

```sql
CREATE USER api_service_account IDENTIFIED BY secret_password;
GRANT CREATE SESSION TO api_service_account;
GRANT SELECT ON SCHEMA my_schema TO api_service_account;
```

## Start API Service

* Download the latest executable from our [GitHub repository](link).

* Start the service running the executable file:
 
```shell
./main
```