default: build

.PHONY: build
build:
	@echo "Generating binary (protoc-gen-openapi) ..."
	@mkdir -p bin
	@go build -o bin/protoc-gen-openapi

.PHONY: fmt
fmt:
	@gofmt -s -w .

PROTO_PATH ?= "internal/converter/testdata/proto"
.PHONY: samples
samples:
	@PATH=./bin:$$PATH protoc --openapi_out=. --proto_path=../../proto/runtime ../../proto/runtime/runtime.proto
