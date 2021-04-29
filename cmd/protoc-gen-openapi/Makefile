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


PROTO_PATH ?= "converter/testdata/proto"
.PHONY: samples
samples:
	@ # Generate a spec for the Micro signup proto:
	@PATH=./bin:$$PATH protoc --openapi_out=. --proto_path=../../proto/signup ../../proto/signup/signup.proto
	@ # Generate a spec for the sample version of our own Micro signup proto (with a different service name):
	@PATH=./bin:$$PATH protoc --openapi_out=service=signuptest:. --proto_path=./converter/testdata/proto ./converter/testdata/proto/signup.proto
