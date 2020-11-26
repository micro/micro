default: generate_spec

generate_spec:
	@protoc --openapi_out=. --proto_path=../proto/auth ../proto/auth/auth.proto
	@protoc --openapi_out=. --proto_path=../proto/config ../proto/config/config.proto
	@protoc --openapi_out=. --proto_path=../proto/signup ../proto/signup/signup.proto
	@protoc --openapi_out=. --proto_path=../proto/store ../proto/store/store.proto
