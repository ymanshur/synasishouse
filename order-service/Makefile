

APP_NAME?=synasishouse-order
APP_IMAGE?=ymanshur/synasishouse:order-latest

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
	-e APP_DB_HOST=host.docker.internal \
	${APP_IMAGE}:latest

POSTGRES_VERSION?=14-alpine
POSTGRES_NAME=postgres${POSTGRES_VERSION}
POSTGRES_USER?=postgres
POSTGRES_PASSWORD?=postgres
POSTGRES_DB?=synansishouse_order
POSTGRES_HOST?=localhost
POSTGRES_PORT?=5432

postgres:
	docker run --name ${POSTGRES_NAME} -p ${POSTGRES_PORT}:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:${POSTGRES_VERSION}

createdb:
	docker exec -it ${POSTGRES_NAME} createdb ${POSTGRES_DB} -U ${POSTGRES_USER}

dropdb:
	docker exec -it ${POSTGRES_NAME} dropdb ${POSTGRES_DB}

psql: clear
	docker exec -it ${POSTGRES_NAME} psql -d ${POSTGRES_DB} -U ${POSTGRES_USER}

DB_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
DB_MIGRATION_PATH?=db/migration

migrate:
	migrate create -ext sql -dir ${DB_MIGRATION_PATH} -seq $(name)

migrateup:
	migrate -path ${DB_MIGRATION_PATH} -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path ${DB_MIGRATION_PATH} -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path ${DB_MIGRATION_PATH} -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path ${DB_MIGRATION_PATH} -database "$(DB_URL)" -verbose down 1

migratedrop:
	migrate -path ${DB_MIGRATION_PATH} -database "$(DB_URL)" -verbose down -all

migratereset: migratedrop migrateup

sqlc:
	sqlc generate -f db/sqlc.yaml

mock:
	mockgen -package mockrepo -destination internal/repo/mock/repo.go github.com/ymanshur/synasishouse/inventory/internal/repo Repo
