NAME=micro
IMAGE_NAME=micro/$(NAME)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_IMPORT=github.com/micro/micro/v3/cmd
CGO_ENABLED=0
BUILD_DATE=$(shell date +%s)
LDFLAGS=-X $(GIT_IMPORT).BuildDate=$(BUILD_DATE) -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT) -X $(GIT_IMPORT).GitTag=$(GIT_TAG)
IMAGE_TAG=$(GIT_TAG)-$(GIT_COMMIT)
PROTO_FLAGS=--go_opt=paths=source_relative --micro_opt=paths=source_relative
PROTO_PATH=$(GOPATH)/src:.
SRC_DIR=$(GOPATH)/src

all: build

vendor:
	go mod vendor

build:
	go build -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" -o $(NAME)

docker:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

.PHONY: proto
proto:
	find proto/ -name '*.proto' -exec protoc --proto_path=$(PROTO_PATH) $(PROTO_FLAGS) --micro_out=$(SRC_DIR) --go_out=plugins=grpc:$(SRC_DIR) {} \;


vet:
	go vet ./...

test: vet
	go test -v ./...

clean:
	rm -rf ./micro

.PHONY: build clean vet test docker
