# synasishouse

Warehousing Microservices with Choreography Pattern

## Design Architecture

<img width="711" height="312" alt="Synasis House-Page-2 drawio" src="https://github.com/user-attachments/assets/34988184-d0b5-4da4-9ea1-d31782379479" />

There are 2 main services to handle the ordering and warehousing process: **Order** and **Inventory** services.

Order acts as a facade of the upstream service that has responsibilities to orchestrate the user's business, including **ordering** and **notification** broker. It is also supposed to maintain **users' sessions** using Redis (which is scalable enough) and route notification messages.

While Inventory serves as the single source of truth for **product and stock management**. It maintains robustness of product availability and consistency of stock, especially for **distributed transactions** by leveraging PostgreSQL ACID properties.

Communication between services using the **gRPC protocol**, and there is a gateway as a **reverse proxy**, while securing the backend from outside exposure.

## API

Open <https://ymanshur.github.io/simplebank/docs/swagger> to see API documentation based on the gRPC Gateway proto definition
