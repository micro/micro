# Micro CLI

This is a CLI for the microservices toolchain [Micro](https://github.com/micro/micro). 

## Getting Started

### Install

```shell
$ go get github.com/micro/micro
```

### Usage
```shell
$ micro
NAME:
   micro - A microservices toolchain

USAGE:
   micro [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   api		Run the micro API
   registry	Query registry
   query	Query a service method using rpc
   health	Query the health of a service
   list		List items in registry
   get		Get item from registry
   sidecar	Run the micro sidecar
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --server_address ":0"	Bind address for the server. 127.0.0.1:8080 [$MICRO_SERVER_ADDRESS]
   --broker "http"		Broker for pub/sub. http, nats, etc [$MICRO_BROKER]
   --broker_address 		Comma-separated list of broker addresses [$MICRO_BROKER_ADDRESS]
   --registry "consul"		Registry for discovery. kubernetes, consul, etc [$MICRO_REGISTRY]
   --registry_address 		Comma-separated list of registry addresses [$MICRO_REGISTRY_ADDRESS]
   --transport "http"		Transport mechanism used; http, rabbitmq, etc [$MICRO_TRANSPORT]
   --transport_address 		Comma-separated list of transport addresses [$MICRO_TRANSPORT_ADDRESS]
   --help, -h			show help
   --version, -v		print the version
```

### List Services
```shell
$ micro list services
go.micro.srv.example
```

### Get Service
```shell
$ micro get service go.micro.srv.example
go.micro.srv.example

go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6	[::]	62421
```

### Query Service
```shell
$ micro query go.micro.srv.example Example.Call '{"name": "John"}'
{
	"msg": "go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello John"
}
```

### Query Service Health
```shell
$ micro health go.micro.sv.example
node		address:port		status
go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6		[::]:62421		ok
```

### Run the API
```shell
$ micro api
I0523 12:23:23.413940   81384 api.go:131] API Rpc handler /rpc
I0523 12:23:23.414238   81384 api.go:143] Listening on [::]:8080
I0523 12:23:23.414272   81384 server.go:113] Starting server go.micro.api id go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
I0523 12:23:23.414355   81384 rpc_server.go:112] Listening on [::]:51938
I0523 12:23:23.414399   81384 server.go:95] Registering node: go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
```

### Run the SideCar
```shell
micro sidecar --server_name=foo --server_address=127.0.0.1:9090 --healthcheck_url=http://127.0.0.1:9090/_status/health
I0523 12:25:36.229536   85658 car.go:184] Registering foo-6ebf29c0-013e-11e5-b55f-68a86d0d36b6
I0523 12:25:36.241680   85658 car.go:188] Starting sidecar healthchecker
```
