# Order Service

This backend application serves to orchestrate order events and process them through an RPC call to the Inventory service, updating stock as needed.

## Table Content

1. [Requirement](#requirement)
2. [Getting Started](#getting-started)
3. [Documentation](#documentation)

## Requirement

User can create an order of multiple products. The application will acquire the ordered products and waiting for settlement or cancelation.

## Getting Started

### Prerequisite

- Go-lang version 1.24.9
- PostgreSQL 14+

### Setup Database

Start database PostgreSQL 14 container service:

```bash
make postgres POSTGRES_VERSION=14
```

Create `synansishouse_order` database:

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

Environment variables allowed in production:

```shell
APP_ENVIRONMENT=development

APP_HTTP_SERVER_ADDR=0.0.0.0:8000

APP_DB_NAME=synansishouse_order
APP_DB_HOST=localhost
APP_DB_PORT=5432
APP_DB_USER=postgres
APP_DB_PASS=postgres
APP_DB_MIGRATION_URL=file://db/migration

APP_GRPC_CLIENT_INVENTORY_HOST=localhost
APP_GRPC_CLIENT_INVENTORY_PORT=9090
APP_GRPC_CLIENT_INVENTORY_MAX_RETRY=3

APP_RABBITMQ_HOST=localhost
APP_RABBITMQ_VHOST=/
APP_RABBITMQ_PORT=5672
APP_RABBITMQ_USER=guest
APP_RABBITMQ_PASS=guest

APP_RABBITMQ_CANCEL_ORDER_QUEUE=order-delayed
APP_RABBITMQ_CANCEL_ORDER_KEY=orders.cancel.delayed
APP_RABBITMQ_CANCEL_ORDER_EXCHANGE=orders.delayed
```

Build the image

```bash
make build
```

Make sure the environment variables are defined when running the following command, update at [Makefile](./Makefile), and the Inventory service run at shared network.

```bash
make run
```

## Documentation

### Data Model

#### Order

| Field | Type | Description | Constraint |
| - | - | - | - |
| id | UUIDv4 | Order internal indetifier | PK |
| order_no | String | Order external identifier | Required, Unique |
| user_id | UUID | User who request order | Required |
| status | String | Status of order | Values: `pending`, `settled`, `canceled` |
| expired_at | Timestamp | When the order expired | |
| updated_at | Timestamp | Last time order was updated | Default: `now()` |
| created_at | Timestamp | Time order was created | Default: `now()` |

#### Order Detail

| Field | Type | Description | Constraint |
| - | - | - | - |
| id | UUIDv4 | Order internal indetifier | PK |
| order_id | UUIDv4 | Order reference indetifier | FK |
| product_code | String | Product external identifier | Required, Unique |
| amount | Number | Amount of product order | Required, Positive |
| updated_at | Timestamp | Last time order detail was updated | Default: `now()` |
| created_at | Timestamp | Time order detail was created | Default: `now()` |

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
    "order_no": "O003",
    "user_id": "6c66959b-4cd1-487c-b010-04dde8616cb6",
    "details": [
        {
            "product_code": "P002",
            "amount": 1
        }
    ]
}
```

#### Success order

```bash
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "code": 200,
  "data": {
    "order_no": "O002",
    "user_id": "6c66959b-4cd1-487c-b010-04dde8616cb6",
    "status": "pending",
    "details": [
      {
        "product_code": "P002",
        "amount": 10
      }
    ]
  },
  "message": "order created successfuly"
}
```

#### Order already exists

```bash
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json; charset=utf-8

{
  "code": 422,
  "message": "order unique constraint violated"
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

### Settle order

```http
POST http://localhost:8000/api/orders/O001/settle HTTP/1.1
Content-Type: application/json
Accept: application/json

{
    "user_id": "6c66959b-4cd1-487c-b010-04dde8616cb6"
}
```

#### Success settle order

```bash
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "code": 200,
  "data": {
    "order_no": "O001",
    "user_id": "6c66959b-4cd1-487c-b010-04dde8616cb6",
    "status": "settled"
  },
  "message": "order settled successfuly"
}
```

#### Settled order already exists

```bash
HTTP/1.1 409 Conflict
Content-Type: application/json; charset=utf-8

{
  "code": 409,
  "message": "settled order already exists"
}

```

### Cancel order

```http
POST http://localhost:8000/api/orders/O003/cancel HTTP/1.1
Content-Type: application/json
Accept: application/json

{
    "user_id": "6c66959b-4cd1-487c-b010-04dde8616cb6"
}
```

#### Success cancel order

```bash
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "code": 200,
  "data": {
    "order_no": "O003",
    "user_id": "6c66959b-4cd1-487c-b010-04dde8616cb6",
    "status": "canceled"
  },
  "message": "order canceled successfuly"
}
```

#### Order already settled

```bash
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json; charset=utf-8

{
  "code": 422,
  "message": "order already settled"
}
```
