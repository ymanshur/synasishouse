# Inventory Service

This backend application serves to manage product stock and provide an RPC endpoint to action warehousing management.

## Table Content

1. [Requirement](#requirement)
2. [Getting Started](#getting-started)
3. [Documentation](#documentation)

## Requirement

### Functional

1. Order service might acquire a stock via a checkout endpoint. The application will **hold stock for later reserve**.
2. If at any time there is a failure for the next process, the hold **stock can be released**.
3. For further development, stock can only be held for a certain period of time, such as 24 hours. After that time, the application will hold back the stock and will notify the Order service or send it as an event to a message broker.

    The **expiry system** that will be applied by scheduling and batching processes. It is carried out so as not to interfere with organic traffic.

### Non Functional

- Exposes gRPC endpoints.
- Cache the product data for total stock in displaying interface using write-through pattern

## Getting Started

### Prerequisite

- Go-lang version 1.24.9
- PostgreSQL 14+
- [Mockery](https://github.com/vektra/mockery) 3+

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

Build the image

```bash
make build
```

Make sure the environment variables are defined when running the following command, update at [Makefile](./Makefile)

```bash
make run
```

## Documentation

### Data Model

#### Product

| Field | Type | Description | Constraint |
| - | - | - | - |
| id | UUIDv4 | Product internal indetifier | PK |
| code | String | Product external identifier | Required, Unique |
| total | Number | Total stock of product | Required, Positive |
| hold | Number | Total hold stock of product | Default: 0 |
| updated_at | Timestamp | Last time product was updated | Default: `now()` |
| created_at | Timestamp | Time product was created | Default: `now()` |

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
    "stocks": [
        {
            "product_code": "P001",
            "amount": 10
        }
    ]
}
EOM
```

### Release Stock

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/ReleaseStock <<EOM
{
    "stocks": [
        {
            "product_code": "P001",
            "amount": 10
        }
    ]
}
EOM
```

### Reserve Stock

```bash
grpcurl -plaintext -d @ localhost:9090 synasishouse.api.Inventory/ReserveStock <<EOM
{
    "stocks": [
        {
            "product_code": "P001",
            "amount": 10
        }
    ]
}
EOM
```
