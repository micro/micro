# micro cli

The **micro cli** is a command line interface for the cloud-native toolkit [micro](https://github.com/micro/micro). 

## Getting Started

### Install

```shell
go get github.com/micro/micro
```

### Usage
```shell
NAME:
   micro - A cloud-native toolkit

USAGE:
   micro [global options] command [command options] [arguments...]
   
VERSION:
   latest
   
COMMANDS:
    api		Run the micro API
    bot		Run the micro bot
    registry	Query registry
    call	Call a service or function
    query	Deprecated: Use call instead
    stream	Create a service or function stream
    health	Query the health of a service
    stats	Query the stats of a service
    list	List items in registry
    register	Register an item in the registry
    deregister	Deregister an item in the registry
    get		Get item from registry
    proxy	Run the micro proxy
    new		Create a new micro service by specifying a directory path relative to your $GOPATH
    run		Run the micro runtime
    web		Run the micro web app
   
GLOBAL OPTIONS:
   --server_name 								Name of the server. go.micro.srv.example [$MICRO_SERVER_NAME]
   --server_version 								Version of the server. 1.1.0 [$MICRO_SERVER_VERSION]
   --server_id 									Id of the server. Auto-generated if not specified [$MICRO_SERVER_ID]
   --server_address 								Bind address for the server. 127.0.0.1:8080 [$MICRO_SERVER_ADDRESS]
   --server_advertise 								Used instead of the server_address when registering with discovery. 127.0.0.1:8080 [$MICRO_SERVER_ADVERTISE]
   --server_metadata [--server_metadata option --server_metadata option]	A list of key-value pairs defining metadata. version=1.0.0 [$MICRO_SERVER_METADATA]
   --broker 									Broker for pub/sub. http, nats, rabbitmq [$MICRO_BROKER]
   --broker_address 								Comma-separated list of broker addresses [$MICRO_BROKER_ADDRESS]
   --registry 									Registry for discovery. memory, consul, etcd, kubernetes [$MICRO_REGISTRY]
   --registry_address 								Comma-separated list of registry addresses [$MICRO_REGISTRY_ADDRESS]
   --selector 									Selector used to pick nodes for querying. random, roundrobin, blacklist [$MICRO_SELECTOR]
   --transport 									Transport mechanism used; http, rabbitmq, nats [$MICRO_TRANSPORT]
   --transport_address 								Comma-separated list of transport addresses [$MICRO_TRANSPORT_ADDRESS]
   --enable_tls									Enable TLS [$MICRO_ENABLE_TLS]
   --tls_cert_file 								TLS Certificate file [$MICRO_TLS_CERT_File]
   --tls_key_file 								TLS Key file [$MICRO_TLS_KEY_File]
   --api_address 								Set the api address e.g 0.0.0.0:8080 [$MICRO_API_ADDRESS]
   --proxy_address 								Proxy requests via the HTTP address specified [$MICRO_PROXY_ADDRESS]
   --sidecar_address 								Set the sidecar address e.g 0.0.0.0:8081 [$MICRO_SIDECAR_ADDRESS]
   --web_address 								Set the web UI address e.g 0.0.0.0:8082 [$MICRO_WEB_ADDRESS]
   --register_ttl "0"								Register TTL in seconds [$MICRO_REGISTER_TTL]
   --register_interval "0"							Register interval in seconds [$MICRO_REGISTER_INTERVAL]
   --api_handler 								Specify the request handler to be used for mapping HTTP requests to services. e.g api, proxy [$MICRO_API_HANDLER]
   --api_namespace 								Set the namespace used by the API e.g. com.example.api [$MICRO_API_NAMESPACE]
   --web_namespace 								Set the namespace used by the Web proxy e.g. com.example.web [$MICRO_WEB_NAMESPACE]
   --api_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_API_CORS]
   --web_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_WEB_CORS]
   --sidecar_cors 								Comma separated whitelist of allowed origins for CORS [$MICRO_SIDECAR_CORS]
   --enable_stats								Enable stats [$MICRO_ENABLE_STATS]
   --help, -h									show help
```

### List Services
```shell
micro list services
```

### Get Service
```shell
micro get service go.micro.srv.example

go.micro.srv.example

go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6	[::]	62421
```

### Call Service
```shell
micro call go.micro.srv.example Example.Call '{"name": "John"}'

{
	"msg": "go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello John"
}
```

### Query Service Health
```shell
micro health go.micro.sv.example

node		address:port		status
go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6		[::]:62421		ok
```

### Register/Deregister with the CLI
```shell
micro register service '{"name": "foo", "version": "bar", "nodes": [{"id": "foo-1", "address": "127.0.0.1", "port": 8080}]}'
```

```shell
micro get service foo

service  foo

version  bar

Id    Address    Port    Metadata
foo-1    127.0.0.1    8080
```

```shell
micro deregister service '{"name": "foo", "version": "bar", "nodes": [{"id": "foo-1", "address": "127.0.0.1", "port": 8080}]}'
```

```shell
micro get service foo

Service not found
```

### Proxy CLI

Proxy remote environments using the `micro proxy`

When developing against remote environments you may not have direct access to service discovery 
which makes it difficult to use the CLI. The `micro proxy` provides a http proxy for such scenarios.

Run the proxy in your remote environment

```
micro proxy
```

Set the env var `MICRO_PROXY_ADDRESS` so the cli knows to use the proxy

```shell
MICRO_PROXY_ADDRESS=staging.micro.mu:8081 micro list services
```
