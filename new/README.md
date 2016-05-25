# Micro New [service]

Micro New is basically a way to generate a boilerplate template. That's it.

## Usage

Create a new service by specifying a directory path relative to your $GOPATH

```
micro new github.com/micro/foo
```

Here it is in action

```
$ micro new github.com/micro/foo
creating service go.micro.srv.foo
creating /Users/asim/checkouts/src/github.com/micro/foo
creating /Users/asim/checkouts/src/github.com/micro/foo/main.go
creating /Users/asim/checkouts/src/github.com/micro/foo/handler
creating /Users/asim/checkouts/src/github.com/micro/foo/handler/example.go
creating /Users/asim/checkouts/src/github.com/micro/foo/subscriber
creating /Users/asim/checkouts/src/github.com/micro/foo/subscriber/example.go
creating /Users/asim/checkouts/src/github.com/micro/foo/proto/example
creating /Users/asim/checkouts/src/github.com/micro/foo/proto/example/example.proto
creating /Users/asim/checkouts/src/github.com/micro/foo/Dockerfile
creating /Users/asim/checkouts/src/github.com/micro/foo/README.md

download protobuf for micro:

go get github.com/micro/protobuf/{proto,protoc-gen-go}

compile the proto file example.proto:

protoc -I/Users/asim/checkouts/src \
	--go_out=plugins=micro:/Users/asim/checkouts/src \
	/Users/asim/checkouts/src/github.com/micro/foo/proto/example/example.proto

```

### Options

Specify more options such as namespace, type, fqdn and alias

```
$ micro new --fqdn com.example.srv.foo github.com/micro/foo
```

### Help

```
NAME:
   micro new - Create a new micro service

USAGE:
   micro new [command options] [arguments...]

OPTIONS:
   --namespace "go.micro"	Namespace for the service e.g com.example
   --type "srv"			Type of service e.g api, srv, web
   --fqdn 			FQDN of service e.g com.example.srv.service (defaults to namespace.type.alias)
   --alias 			Alias is the short name used as part of combined name if specified
```
