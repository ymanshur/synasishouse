# synasishouse

Warehousing Microservices with Choreography Pattern

## Design Architecture

<img width="681" height="311" alt="SynansisHouse drawio" src="https://github.com/user-attachments/assets/ac4db1c7-4430-4267-a105-72177609581e" />

## API

Open <https://ymanshur.github.io/simplebank/docs/swagger> to see API documentation based on the gRPC Gateway proto definition

### Checkout an order

```http
POST http://0.0.0.0:8000/checkout HTTP/1.1
Content-Type: application/json
Accept: application/json

{
    "code": "NO01",
    "amount": 1
}
```
