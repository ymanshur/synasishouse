# synasishouse

Warehousing Microservices with Choreography Pattern

***

There are 2 main services to handle the ordering and warehousing process: **Order** and **Inventory** services.

The Order acts as a facade for the upstream service, which has responsibilities to orchestrate the user's business, including ordering and notification brokering.

While Inventory serves as the single source of truth for **product and stock management**. It maintains robustness of product availability and consistency of stock, especially for **distributed transactions** by leveraging **PostgreSQL** ACID properties.

Communication between services using the **gRPC protocol**, and there is a gateway as a **reverse proxy**, while securing the backend from outside exposure.

Each reserved stock is treated like a **transaction that has an expiry**. Every time a transaction expires, the Inventory service will publish the event and be consumed by the Order service, or at the time when the transaction will be reserved, then Inventory will send an invalid transaction error according to their expiry.

## Design Architecture

### Simplicity of Sake

<img width="646" height="311" alt="SynasisHouse-BAU" src="https://github.com/user-attachments/assets/ac1eb332-d9e7-45af-b6a9-bc8dddb9b99d" />

### Choreography Pattern

<img width="733" height="682" alt="SynasisHouse-BEYOND" src="https://github.com/user-attachments/assets/071cc190-62d3-4043-88e4-c84e3dc22c2d" />

### User and Product Management

It should be considered to provide the **User** and **Product** services independently. User service stores user profiles to PostgreSQL and their sessions in **Redis** (which is scalable enough) and routes notifications to the **Notification** service. Product service is the place where Admin ships the product or product stocks to the Inventory service.

### Observability

Elasticsearch and Filebeat

## Data Model

<img width="720" height="682" alt="Synasis House" src="https://github.com/user-attachments/assets/48efdb90-bd79-47d2-b8a1-b40c17050f0a" />

More details in: <https://dbdocs.io/ymanshur/SynasisHouse>

### Notable Decisions

 - Only exposes shared references, not internal identifiers, i.e., `order_no` and `product_code`. As preparation for beyond service separation.
 - Reference inventory and transaction are optional; consider decouple each other to minimize lock contention between transactional behaviour and expiry processing, as long as `product_code` never changes.

## Getting Started

```bash
make compose
```

### Container Registry (hub.docker.com)

Pull the images from <https://hub.docker.com/repository/docker/ymanshur/synasishouse>

#### Inventory Service

```bash
docker pull ymanshur/synasishouse:inventory-latest
```

#### Order Service

```bash
docker pull ymanshur/synasishouse:order-latest
```

## API

Open <https://ymanshur.github.io/synasishouse/docs/swagger/> to see API documentation based on the gRPC Gateway proto annotation.
