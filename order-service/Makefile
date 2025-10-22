

APP_NAME?=order
APP_IMAGE?=ymanshur/synasishouse/${APP_NAME}

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

# run with necessary environment variables
run:
	docker run --rm \
	--name ${APP_NAME} \
	-p 8000:8000 \
	-e APP_ENVIRONMENT=development \
	${APP_IMAGE}:latest

