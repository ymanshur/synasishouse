

APP_NAME?=inventory
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
	-p 9090:9090 \
	-e APP_DB_SOURCE="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@host.docker.internal:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
	${APP_IMAGE}:latest

POSTGRES_VERSION?=14-alpine
POSTGRES_NAME=postgres${POSTGRES_VERSION}
POSTGRES_USER?=postgres
POSTGRES_PASSWORD?=postgres
POSTGRES_DB?=inventory
POSTGRES_HOST?=localhost
POSTGRES_PORT?=5432

postgres:
	docker run --name ${POSTGRES_NAME} -p ${POSTGRES_PORT}:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:${POSTGRES_VERSION}

createdb:
	docker exec -it ${POSTGRES_NAME} createdb ${POSTGRES_DB} -U ${POSTGRES_USER}

dropdb:
	docker exec -it ${POSTGRES_NAME} dropdb ${POSTGRES_DB}

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

mocks:
	mockgen -package mockrepo -destination internal/repo/mock/repo.go github.com/ymanshur/synasishouse/inventory/internal/repo Repo
