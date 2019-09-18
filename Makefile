NAME=micro
IMAGE_NAME=micro/$(NAME)
TAG=$(shell git describe --abbrev=0 --tags)
CGO_ENABLED=0

all: build

vendor:
	go mod vendor

build:
	go get
	go build -a -installsuffix cgo -ldflags '-w' -o $(NAME) ./*.go

docker:
	docker build -t $(IMAGE_NAME):$(TAG) .
	docker tag $(IMAGE_NAME):$(TAG) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(TAG)
	docker push $(IMAGE_NAME):latest

vet:
	go vet ./...

test: vet
	go test -v ./...

clean:
	rm -rf ./micro

.PHONY: build clean vet test docker
