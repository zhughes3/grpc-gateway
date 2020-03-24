# GRPC Server w/ HTTP Gateway

### Init Postgres DB
```text
# Start psql as superadmin user
psql -U postgres

CREATE DATABASE <DATABASE_NAME>;

CREATE USER <USERNAME> WITH PASSWORD '<PASSWORD>';
```

### Create configuration file
To work properly, a `config.env` file needs to be placed in the root directory. At a minimum, this file needs the following variables set:
```.env
RPC_PORT=<RPC_PORT>
HTTP_PORT=<HTTP_PORT>
DB_USER=<DB_USERNAME>
DB_PASSWORD=<DB_PASSWORD>
DB_PORT=<DB_PORT>
DB_NAME=<NAME OF DATABASE>
```

### Generate Golang protocol buffer code
```bash
./gen-proto.sh
```

### Start Application
```bash
go run .
```