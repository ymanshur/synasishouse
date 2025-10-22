# Order Service

This backend application serves to orchestrate order events and process them to order service though RPC call: **Checkout**, **Reserve**, and **Release** to validatate stock.

## Table Content

1. [Requirement](#requirement)
2. [Getting Started](#getting-started)
3. [Documentation](#documentation)

## Requirement

### Functional

### Non Functional

## Getting Started

### Prerequisite

- Go-lang version 1.24.9
- PostgreSQL 14

### Setup Database

Start database PostgreSQL 14 container service:

```bash
make postgres POSTGRES_VERSION=14
```

Create `order` database:

```bash
make createdb
```

Create database migration under database/migration:

```bash
make migrate name=init_schema
```

### Run the backend on the local machine

Copy the configuration file under the config directory and run:

```bash
cp config/app.env.example config/app.env
```

```bash
make server
```

It will run at <http://0.0.0.0:8000> as the default

Test the backend:

```bash
make test
```

### Run the backend and the database in a Docker container

Environment variables allowed for production service:

```shell
ENVIRONMENT=development

HTTP_SERVER_ADDRESS=0.0.0.0:8000

GRPC_CLIENT_HOST_INVENTORY=localhost
GRPC_CLIENT_PORT_INVENTORY=9090
```

Make sure the environment variables are defined when running the following command, update at [Makefile](./Makefile)

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
