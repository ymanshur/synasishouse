.PHONY: proto
proto:
	rm -f proto/*.go
	rm -f docs/swagger/*.swagger.json
	protoc \
	--proto_path=proto \
	--go_out=proto --go_opt=paths=source_relative \
    --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,json_names_for_fields=false \
    --experimental_allow_proto3_optional \
    proto/*.proto

compose:
	mkdir -p rabbitmq/data
	mkdir -p rabbitmq/log
	mkdir -p storage/inventory/pg_data
	mkdir -p storage/order/pg_data
	docker compose up

RABBITMQ_VERSION?=3-management-alpine
RABBITMQ_NAME=rabbitmq${RABBITMQ_VERSION}
RABBITMQ_USER?=guest
RABBITMQ_PASS?=guest

.PHONY: rabbitmq
rabbitmq:
	docker run -d --name ${RABBITMQ_NAME} \
		-p 5672:5672 -p 15672:15672 \
		-v ./rabbitmq/data/:/var/lib/rabbitmq/ \
      	-v ./rabbitmq/log:/var/log/rabbitmq/ \
      	-v ./rabbitmq/enabled_plugins:/etc/rabbitmq/enabled_plugins:rw \
      	-v ./rabbitmq/plugins:/usr/lib/rabbitmq/plugins \
		-e RABBITMQ_DEFAULT_USER=${RABBITMQ_USER} \
		-e RABBITMQ_DEFAULT_PASS=${RABBITMQ_PASS} \
		-e RABBITMQ_PLUGINS_DIR=/opt/rabbitmq/plugins:/usr/lib/rabbitmq/plugins \
		rabbitmq:${RABBITMQ_VERSION}

rabbitmq-status:
	docker exec ${RABBITMQ_NAME} rabbitmqctl status
