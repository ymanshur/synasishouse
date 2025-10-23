

APP_NAME?=synasishouse-order
APP_IMAGE?=ymanshur/synasishouse/order

clear:
	clear

server: clear
	go run main.go

.PHONY: test
test:
	go test -v -cover ./...

modsync:
	go get -u
	go mod tidy
	go mod vendor

build: modsync
	docker build -t ${APP_IMAGE}:latest .

NETWORK?=synasishouse-net

network:
	docker network create ${NETWORK}

# run with necessary environment variables
run:
	docker run --rm \
	--name ${APP_NAME} \
	--network ${NETWORK} \
	-p 8000:8000 \
	-e APP_ENVIRONMENT=development \
	-e APP_GRPC_CLIENT_INVENTORY_HOST=synasishouse-inventory \
	${APP_IMAGE}:latest

