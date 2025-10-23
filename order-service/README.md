# Order Service

This backend application serves to orchestrate order events and process them to Inventory service though RPC call to action stock.

## Table Content

1. [Requirement](#requirement)
2. [Getting Started](#getting-started)
3. [Documentation](#documentation)

## Requirement

<!-- ### Functional

### Non Functional -->

## Getting Started

### Prerequisite

- Go-lang version 1.24.9

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

Environment variables allowed in production:

```shell
APP_ENVIRONMENT=development

APP_HTTP_SERVER_ADDR=0.0.0.0:8000

APP_GRPC_CLIENT_INVENTORY_HOST=localhost
APP_GRPC_CLIENT_INVENTORY_PORT=9090
```

Build the image with vendor mode (it's mandatory)

```bash
make build
```

Make sure the environment variables are defined when running the following command, update at [Makefile](./Makefile), and the Inventory service run at shared network.

```bash
make run
```

## Documentation

<!-- ### Data Model -->

### API

#### Health Check

```http
GET http://localhost:8000/api/health HTTP/1.1
```

### Create order

```http
POST http://localhost:8000/api/orders HTTP/1.1
Content-Type: application/json
Accept: application/json

{
    "code": "P002",
    "amount": 1
}
```

#### Product not found

```bash
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8

{
  "code": 404,
  "message": "product not found"
}
```

#### Stock is unavailable

```bash
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json; charset=utf-8

{
  "code": 422,
  "message": "stock is unavailable"
}
```

#### Stock is available

```bash
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "code": 200,
  "message": "stock is available"
}
```
