# Inventory Service

This backend application serves to manage product stock and provide an RPC endpoint to action warehousing management.

## Table Content

1. [Requirement](#requirement)
2. [Getting Started](#getting-started)
3. [Documentation](#documentation)

## Requirement

### Functional

1. Other parties might acquire a stock via a checkout endpoint. The application will **hold reserved stock for later confirmation**.
2. If at any time there is a failure for the next process, the reserved **stock can be released**.
3. For further development, stock can only be held for a certain period of time, such as 24 hours. After that time, the application will release the stock and will notify the *Order* application or send it in the form of an event to a message broker.

    The **expiry system** that will be applied;
    1. Synchronously, every time a product is checked out or released, along with other stock.
    2. Scheduling and batching processes are carried out so as not to interfere with organic traffic.

### Non Functional

- Exposes gRPC endpoints.

## Getting Started

### Prerequisite

- Go-lang version 1.24.9
- PostgreSQL 14

### Setup Database

Start database PostgreSQL 14 container service:

```bash
make postgres POSTGRES_VERSION=14
```

Create `synansishouse_inventory` database:

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

It will run at <http://0.0.0.0:9090> as the default

Test the business logic:

```bash
make test
```

### Run the backend and the database in a Docker container

Environment variables allowed in production:

```shell
APP_GRPC_SERVER_ADDR=0.0.0.0:9090

APP_DB_NAME=synansishouse_inventory
APP_DB_HOST=localhost
APP_DB_PORT=5432
APP_DB_USER=postgres
APP_DB_PASS=postgres
APP_DB_MIGRATION_URL=file://db/migration
```

Build the image with vendor mode (it's mandatory)

```bash
make build
```

Create a network to communicate with Order service

```bash
make network
```

Make sure the environment variables are defined when running the following command, update at [Makefile](./Makefile)

```bash
make run
```

## Documentation

### Data Model

<img width="377" height="344" alt="Synasis House" src="https://github.com/user-attachments/assets/2d8e3f63-39d0-4807-9d75-03e7f1c28b7a" />

#### Product

| Field | Type | Description | Constraint |
| - | - | - | - |
| id | UUIDv4 | Product internal indetifier | PK |
| code | String | Product external identifier | Required, Unique |
| total | Number | Total stock of product | Required, Positive |
| reserved | Number | Total reserved stock of product | Default: 0 |

## API

Install [grpcurl](https://github.com/fullstorydev/grpcurl) for local test (optional)

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### Health Check

```bash
grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check
```

### Create Product

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/CreateProduct <<EOM
{
    "code": "P002",
    "total": 100
}
EOM
```

### Get Product

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/GetProduct <<EOM
{
    "id": "6e4418dd-81d7-469f-8f36-b29f8b741b8c"
}
EOM
```

### Update Product

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/UpdateProduct <<EOM
{
    "code": "P003",
    "id": "6e4418dd-81d7-469f-8f36-b29f8b741b8c"
}
EOM
```

### Delete Product

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/DeleteProduct <<EOM
{
    "id": "6e4418dd-81d7-469f-8f36-b29f8b741b8c"
}
EOM
```

### Check Stock

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/CheckStock <<EOM
{
    "code": "P001",
    "amount": 10
}
EOM
```

### Reserve Stock

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/ReserveStock <<EOM
{
    "code": "P001",
    "amount": 10
}
EOM
```

### Release Stock

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/ReleaseStock <<EOM
{
    "code": "P001",
    "amount": 10
}
EOM
```
