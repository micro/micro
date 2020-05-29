package template

var (
	Makefile = `
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto
{{if ne .Type "web"}}
.PHONY: proto
proto:
    {{if eq .UseGoPath true}}
	protoc --proto_path=${GOPATH}/src:. --micro_out=${MODIFY}:. --go_out=${MODIFY}. proto/{{.Alias}}/{{.Alias}}.proto
    {{else}}
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/{{.Alias}}/{{.Alias}}.proto
    {{end}}

.PHONY: build
build: proto
{{else}}
.PHONY: build
build:
{{end}}
	go build -o {{.Alias}}-{{.Type}} *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t {{.Alias}}-{{.Type}}:latest
`

	GenerateFile = `package main
//go:generate make proto
`
)
