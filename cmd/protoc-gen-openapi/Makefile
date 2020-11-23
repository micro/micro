default: build

.PHONY: build
build:
	@echo "Generating binary (protoc-gen-openapi) ..."
	@mkdir -p bin
	@go build -o bin/protoc-gen-openapi

.PHONY: fmt
fmt:
	@gofmt -s -w .

.PHONY: install
install:
	@GO111MODULE=on go install


PROTO_PATH ?= "internal/converter/testdata/proto"
.PHONY: samples
samples:
	# @PATH=./bin:$$PATH protoc --openapi_out=. --proto_path=../../proto/runtime ../../proto/runtime/runtime.proto
	# @PATH=./bin:$$PATH protoc --openapi_out=service=event:. --proto_path=../../proto/events ../../proto/events/events.proto
	@PATH=./bin:$$PATH protoc --openapi_out=. --proto_path=../../proto/signup ../../proto/signup/signup.proto