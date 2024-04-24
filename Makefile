NAME = micro
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_TAG = $(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_IMPORT = micro.dev/v4/cmd
BUILD_DATE = $(shell date +%s)
LDFLAGS = -X $(GIT_IMPORT).BuildDate=$(BUILD_DATE) -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT) -X $(GIT_IMPORT).GitTag=$(GIT_TAG)

DOCKER_BUILD = docker buildx build
DOCKER_BUILD_ARGS = --platform linux/amd64 --platform linux/arm64
DOCKER_IMAGE_NAME = micro/$(NAME)
DOCKER_IMAGE_TAG = --tag $(DOCKER_IMAGE_NAME):$(GIT_TAG)-$(GIT_COMMIT) --tag $(DOCKER_IMAGE_NAME):latest

PROTO_FILES = $(wildcard proto/**/*.proto)
PROTO_GO_MICRO = $(PROTO_FILES:.proto=.pb.go) $(PROTO_FILES:.proto=.pb.micro.go)


.DEFAULT_GOAL := $(NAME)

.PHONY: tidy
tidy:
	go mod tidy

$(NAME):
	CGO_ENABLED=0 go build -ldflags "-s -w ${LDFLAGS}" -o $(NAME) cmd/micro/main.go

.PHONY: docker
docker:
	$(DOCKER_BUILD) $(DOCKER_BUILD_ARGS) $(DOCKER_IMAGE_TAGS) --push .

.PHONY: proto
proto: $(PROTO_GO_MICRO)

%.pb.micro.go %.pb.go: %.proto clean
	protoc --proto_path=. --micro_out=. --go_out=. $<

.PHONY: test
vet:
	go vet ./...

.PHONY: vet
test: vet
	go test -v -race ./...

.PHONY: clean
clean:
	rm -f $(NAME) $(PROTO_GO_MICRO)

.PHONY: gorelease-dry-run
gorelease-dry-run:
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-v $(CURDIR):/$(NAME) \
		-w /$(NAME) \
		ghcr.io/goreleaser/goreleaser-cross:v1.20.6 \
		--clean --skip-validate --skip-publish
