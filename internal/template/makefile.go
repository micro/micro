package template

var (
	Makefile = `
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker

{{if ne .Type "web"}}
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/example/example.proto

build: proto
{{else}}
build:
{{end}}
	go build -o {{.Alias}}-{{.Type}} main.go plugin.go

test:
	go test -v ./... -cover

docker:
	docker build . -t {{.Alias}}-{{.Type}}:latest
`
)
