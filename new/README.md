# micro new [service]

The **micro new** command is a quick way to generate boilerplate templates for micro services.

## Usage

Create a new service by specifying a directory path relative to your $GOPATH

```
micro new github.com/micro/example
```

Here it is in action

```
micro new github.com/micro/example

Creating service go.micro.srv.example in /home/go/src/github.com/micro/example 

.
├── main.go
├── handler
│   └── example.go
├── subscriber
│   └── example.go
├── proto/example
│   └── example.proto
├── Dockerfile
└── README.md


download protobuf for micro:

brew install protobuf
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u github.com/micro/protoc-gen-micro

compile the proto file example.proto:

cd /home/go/src/github.com/micro/example
protoc --proto_path=. --go_out=. --micro_out=. proto/example/example.proto
```

### Options

Specify more options such as namespace, type, fqdn and alias

```
micro new --fqdn io.foobar.srv.example github.com/micro/example
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
