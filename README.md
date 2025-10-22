# synasishouse

Warehousing Microservices with Choreography Pattern

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
