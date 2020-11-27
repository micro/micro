default: create_dirs generate_spec generate_docs

create_dirs:
	@mkdir -p auth broker config events registry runtime signup store

generate_spec:
	@protoc --openapi_out=auth --proto_path=../proto/auth ../proto/auth/auth.proto
	@protoc --openapi_out=broker --proto_path=../proto/broker ../proto/broker/broker.proto
	@protoc --openapi_out=config --proto_path=../proto/config ../proto/config/config.proto
	@protoc --openapi_out=events --proto_path=../proto/events ../proto/events/events.proto
	# @protoc --openapi_out=registry --proto_path=../proto/registry ../proto/registry/registry.proto
	@protoc --openapi_out=runtime --proto_path=../proto/runtime ../proto/runtime/runtime.proto
	@protoc --openapi_out=signup --proto_path=../proto/signup ../proto/signup/signup.proto
	@protoc --openapi_out=store --proto_path=../proto/store ../proto/store/store.proto

generate_docs:
	@redoc-cli bundle auth/api-auth.json
