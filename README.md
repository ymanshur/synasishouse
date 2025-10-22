# synasishouse

Warehousing Microservices with Choreography Pattern

## Design Architecture

<img width="711" height="331" alt="Synasis House-Page-2 drawio" src="https://github.com/user-attachments/assets/72e0164b-eb76-4c86-91a7-165d5ea8b326" />

There are 2 main services to handle the ordering and warehousing process: **Order** and **Inventory** services.

The Order acts as a facade for the upstream service, which has responsibilities to orchestrate the user's business, including ordering and notification brokering. It is also designed to maintain users' data and their sessions using **Redis** (which is scalable enough) and route notification messages.

While Inventory serves as the single source of truth for **product and stock management**. It maintains robustness of product availability and consistency of stock, especially for **distributed transactions** by leveraging **PostgreSQL** ACID properties.

Communication between services using the **gRPC protocol**, and there is a gateway as a **reverse proxy**, while securing the backend from outside exposure.

<!-- ### Observability -->

## API

Open <https://ymanshur.github.io/synasishouse/docs/swagger/> to see API documentation based on the gRPC Gateway proto annotation.
