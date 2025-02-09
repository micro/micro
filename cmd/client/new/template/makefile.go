package template

var (
	Makefile = `GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOBIN := $(shell go env GOPATH)/bin
NAME := $(shell  basename $(CURDIR))
ifeq ($(GOOS),windows)
	TARGET := $(NAME).exe
else
	TARGET := $(NAME)
endif
PROTO_FILES := $(wildcard proto/*.proto)
PROTO_GO := $(PROTO_FILES:%.proto=%.pb.go)
PROTO_GO_MICRO := $(PROTO_FILES:%.proto=%.pb.micro.go)

.DEFAULT_GOAL := $(TARGET)

# Build binary
$(TARGET): proto/$(NAME).pb.micro.go
	go build -o $(TARGET) *.go

.PHONY: clean
clean:
	rm -f $(TARGET) $(PROTO_GO) $(PROTO_GO_MICRO)

include dep-install.mk

.PHONY: proto
proto: proto/$(NAME).pb.micro.go

%.pb.go %.pb.micro.go: %.proto
	protoc --proto_path=. --go_out=:. --micro_out=. $<

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t $(NAME):latest
`
)
