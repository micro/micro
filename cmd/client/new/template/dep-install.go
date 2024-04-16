package template

var (
	DepInstall = `PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_MICRO := $(GOBIN)/protoc-gen-micro
PROTOC = $(shell which protoc || echo "$(GOBIN)/protoc")

# Generate protoc latest binary url
PROTOC_RELEASE_BIN := protoc
ifeq ($(GOOS), linux)
	ifeq ($(GOARCH), amd64)
		PROTOC_RELEASE = linux-x86_64
	else ifeq ($(GOARCH), arm64)
		PROTOC_RELEASE = linux-aarch_64
	else ifeq ($(GOARCH), 386)
		PROTOC_RELEASE = linux-x86_32
	else
		echo >&2 "only support amd64, arm64 or 386 GOARCH" && exit 1
	endif
else ifeq ($(GOOS), darwin)
	ifeq ($(GOARCH), amd64)
		PROTOC_RELEASE = osx-x86_64
	else ifeq ($(GOARCH), arm64)
		PROTOC_RELEASE = osx-aarch_64
	else
		echo >&2 "only support amd64 or arm64 GOARCH" && exit 1
	endif
else ifeq ($(GOOS), windows)
	ifeq ($(GOARCH), amd64)
		PROTOC_RELEASE = win64
	else ifeq ($(GOARCH), 386)
		PROTOC_RELEASE = win32
	else
		echo >&2 "only support amd64 or 386 GOARCH" && exit 1
	endif
PROTOC_RELEASE_BIN = protoc.exe
endif
PROTOC_LATEST = https://github.com/protocolbuffers/protobuf/releases/latest
PROTOC_VERSION = $(shell curl -Ls -o /dev/null -w %{url_effective} $(PROTOC_LATEST) | sed 's|.*/v||')
PROTOC_RELEASE_URL = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_RELEASE).zip

# Install dependencies
.PHONY: init
init: $(PROTOC) $(PROTOC_GEN_GO) $(PROTOC_GEN_MICRO)

# Install protoc latest binary
# from release https://github.com/protocolbuffers/protobuf/releases
# to $(GOBIN)
$(PROTOC):
	curl -LSs $(PROTOC_RELEASE_URL) -o protoc.zip
	unzip -qqj protoc.zip "bin/$(PROTOC_RELEASE_BIN)" -d "$(GOBIN)" && rm protoc.zip

# Install protoc-gen-go
$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install protoc-gen-micro
$(PROTOC_GEN_MICRO):
	go install micro.dev/v4/cmd/protoc-gen-micro@master
`
)
