# Inventory Service

## Table Content

1. [Requirement](#requirement)
2. [Getting Started](#getting-started)
3. [Documentation](#documentation)

## Requirement

## Getting Started

### Prerequisite

- Go-lang version 1.24.9
- PostgreSQL 14

### Setup Database

Start database PostgreSQL 14 container service:

```bash
make postgres POSTGRES_VERSION=14
```

Create `simplewallet` database:

```bash
make createdb
```

Create database migration under database/migration:

```bash
make migrate name=init_schema
```

### Run the backend on the local machine

Copy the configuration file under config directory and run:

```bash
cp config/app.yaml.dist config/app.yaml
```

```bash
make server
```

It will run at <http://0.0.0.0:8000> as default

Test the backend:

```bash
make test
```

### Run the backend and the database in Docker containers

Environment variables allowed for production service:

```shell
APP_NAME=inventory
APP_ENV=production
APP_DEBUG=false
DB_NAME=inventory
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
```

The following command will create PostgreSQL 14 and bind the volume data in [tmp](tmp) directory and run the [docker-compose.yaml](deployment/docker-compose.yaml) file after build the backend image.

```bash
make containers
```

Alternatively, if you have already have PostgreSQL service, just run the following command to create only the backend container

```bash
make run
```

## Documentation

### Data Model

### API

#### Health Check

```http
GET http://0.0.0.0:8000/health HTTP/1.1
```
