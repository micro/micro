default: generate_spec

generate_spec:
	@protoc --openapi_out=. --proto_path=../proto/auth ../proto/auth/auth.proto
	@protoc --openapi_out=. --proto_path=../proto/broker ../proto/broker/broker.proto
	@protoc --openapi_out=. --proto_path=../proto/config ../proto/config/config.proto
	@protoc --openapi_out=. --proto_path=../proto/events ../proto/events/events.proto
	# @protoc --openapi_out=. --proto_path=../proto/registry ../proto/registry/registry.proto
	@protoc --openapi_out=. --proto_path=../proto/runtime ../proto/runtime/runtime.proto
	@protoc --openapi_out=. --proto_path=../proto/signup ../proto/signup/signup.proto
	@protoc --openapi_out=. --proto_path=../proto/store ../proto/store/store.proto
