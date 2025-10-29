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
	mkdir -p storage/inventory/pg_data
	mkdir -p storage/order/pg_data
	docker compose up
