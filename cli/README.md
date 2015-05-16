# Micro CLI

This is a CLI for the microservices toolchain [Micro](https://github.com/myodc/micro). 

## Getting Started

### Install

```shell
$ go get github.com/myodc/micro
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

AUTHOR(S): 
   
COMMANDS:
   api		Run the micro API
   registry	Query registry
   store	Query store
   list		List items in registry
   get		get item from registry
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

### List Services
```shell
$ micro list services
go.micro.service.template
```

### Get Service
```shell
$ micro get service go.micro.service.template
go.micro.service.template

go.micro.service.template-c5718d29-da2a-11e4-be11-68a86d0d36b6	[::]	60728
```

### Query Service
```shell
$ micro query go.micro.service.template Example.Call '{"name": "John"}'
{
	"msg": "go.micro.service.template-5c3b2801-fc1b-11e4-9f62-68a86d0d36b6: Hello John"
}
```

### Get Item from Store
```shell
$ micro store get foo
bar
```

### Run the API
```shell
$ micro api
I0407 23:14:28.347041   22506 rpc_server.go:156] Rpc handler /_rpc
I0407 23:14:28.347372   22506 api.go:97] API Rpc handler /rpc
I0407 23:14:28.347447   22506 api.go:108] Listening on [::]:8080
I0407 23:14:28.347481   22506 server.go:90] Starting server go.micro.api id go.micro.api-75184b9f-dd73-11e4-937f-68a86d0d36b6
I0407 23:14:28.347570   22506 rpc_server.go:187] Listening on [::]:54120
I0407 23:14:28.347609   22506 server.go:76] Registering go.micro.api-75184b9f-dd73-11e4-937f-68a86d0d36b6
```

### Run the SideCar
```shell
micro sidecar --server_name=foo --server_address=127.0.0.1:9090 --healthcheck_url=http://127.0.0.1:9090/_status/health
I0409 21:25:12.817687   27939 sic.go:182] Registering foo-86857602-def6-11e4-99f5-68a86d0d36b6
I0409 21:25:12.844155   27939 sic.go:186] Starting sidecar healthchecker
```
