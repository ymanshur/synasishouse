# Inventory Service

This backend application serves to manage product stock and provide an RPC endpoint to action warehousing management, including: **Checkout**, **Reserve**, and **Release** stock.

## Table Content

1. [Requirement](#requirement)
2. [Getting Started](#getting-started)
3. [Documentation](#documentation)

## Requirement

### Functional

1. Other parties might acquire a stock via a checkout endpoint. The application will **hold reserved stock for later confirmation**.
2. If at any time there is a failure for the next process, the reserved stock can be released again.
3. For further development, stock can only be held for a certain period of time, such as 24 hours. After that time, the application will release the stock and will notify the *Order* application or send it in the form of an event to a message broker.

### Non Functional

- Communication is done using the **gRPC protocol** and Pub/Sub messages if needed.
- The expiry system that will be applied;
    1. Synchronously, every time a product is checked out or released, along with other stock.
    2. Scheduling and batching processes are carried out so as not to interfere with organic traffic.

## Getting Started

### Prerequisite

- Go-lang version 1.24.9
- PostgreSQL 14

### Setup Database

Start database PostgreSQL 14 container service:

```bash
make postgres POSTGRES_VERSION=14
```

Create `inventory` database:

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

Test the backend:

```bash
make test
```

### Run the backend and the database in a Docker container

Environment variables allowed for production service:

```shell
DB_SOURCE=postgresql://postgres:postgres@localhost:5432/inventory?sslmode=disable
DB_MIGRATION_URL=file://db/migration

GRPC_SERVER_ADDRESS=0.0.0.0:9090
```

Make sure the environment variables are defined when running the following command, update at [Makefile](./Makefile)

```bash
make run
```

## Documentation

### Data Model

<img width="377" height="344" alt="Synasis House" src="https://github.com/user-attachments/assets/2d8e3f63-39d0-4807-9d75-03e7f1c28b7a" />
