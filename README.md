# PDSM - PDS Customer Service

Initialize project for PDS Customer Service

## Prerequisites

1. Go 1.16
2. PostgreSQL 11
3. UNIX Shell
   > Use `wsl` in Windows 10
4. Git
5. Make
6. Docker CE (Optional)

## Set-up

1. Configure Project

   ```sh
   # Run scripts to make env from .env-example and grant permission
   make setup
   
   # Run scripts to get all app dependencies that needed.
   make configure
   
   # Run scripts to check all prerequisites for development is available
   make doctor
   ```

2. Configure project. see [Configuration Section](#Configuration) for details:

3. Init Database

   Once `.env` has been configured, initiate database:
   ```bash
   # Create database if not exists
   make db
   
   # Upgrade database to next version
   make db-up
   ```

## Configuration

PDS service are configurable from `.env` file

Below are available configuration for the project:

### Run Development

```sh
# Run Service
make serve
```

|                      | Key                        | Description                                 | Required                                          | Value                                                                                      |
| -------------------- | -------------------------- | ------------------------------------------- | ------------------------------------------------- |:------------------------------------------------------------------------------------------ |
| **Common**           | `DEBUG`                    | Debug Mode                                  |                                                   | Boolean. Default: `false`                                                                  | 
| **API - Server**     | `PORT`                     | Server listen port                          | **✓**                                             | `1024-65535`                                                                               |
|                      | `SERVER_HOSTNAME`          | Resolved hostname                           |                                                   | String, Default: `localhost`                                                               |
|                      | `SERVER_BASE_PATH`         | Resolved base path                          |                                                   | String, Default: `/`                                                                       |
|                      | `SERVER_LISTEN_SECURE`     | Listen in secure mode for base url resolver |                                                   | Boolean, Default: `false`                                                                  |
|                      | `SERVER_TRUST_PROXY`       | Show debug responses                        | **✓** (Required for Deployment via Reverse Proxy) | String Array of IP Address, separated by comma. Default: `[]`. Set to `["*"]` to allow all |
|                      | `SERVER_HTTP_BASE_URL`     | Override value for resolved HTTP Base URL   |                                                   | URL                                                                                        |
|                      |                            |                                             |                                                   |                                                                                            |
|  **API - Client**    | `CLIENT_ID`                |  Client credential for Web App              | **✓**                                             | String                                                                                     |
|                      | `CLIENT_SECRET`            |  Client credential for Web App              | **✓**                                             | String                                                                                     |
|                      | `CORS_ENABLED`             |  Allowed CORS Origin                        | **✓**                                             | Comma-separated URL. Use value `'*'` to allow all origin                                   |
|                      |                            |                                             |                                                   |                                                                                            |
| **Data Sources**     | `DB_DRIVER`                | Database driver                             | **✓**                                             | `postgres`                                                                                 |
|                      | `DB_HOST`                  | Postgres server host                        | **✓**                                             | String                                                                                     |
|                      | `DB_PORT`                  | Postgres server port                        | **✓**                                             | `1024-65535`                                                                               |
|                      | `DB_USER`                  | Postgres server username                    | **✓**                                             | String                                                                                     |
|                      | `DB_PASS`                  | Postgres server password                    | **✓**                                             | String                                                                                     |
|                      | `DB_NAME`                  | Postgres server database                    | **✓**                                             | String                                                                                     |

## Contributors

- Saggaf Arsyad <saggaf@nusantarabetastudio.com>
- Dinan <dinan@nusantarabetastudio.com>