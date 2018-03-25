# micro cli

The **micro cli** is a command line interface for the cloud-native toolkit [micro](https://github.com/micro/micro). 

## Getting Started

- [Install](#install)
- [List Services](#list-services)
- [Get Service](#get-service)
- [Call Service](#call-service)
- [Service Health](#service-health)
- [Proxy Remote Environment](#proxy-remote-env)

## Install

```shell
go get github.com/micro/micro
```

## Example Usage

### List Services

```shell
micro list services
```

### Get Service

```shell
micro get service go.micro.srv.example
```

Output

```
go.micro.srv.example

go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6	[::]	62421
```

### Call Service

```shell
micro call go.micro.srv.example Example.Call '{"name": "John"}'
```

Output
```
{
	"msg": "go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello John"
}
```

### Service Health

```shell
micro health go.micro.sv.example
```

Output

```
node		address:port		status
go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6		[::]:62421		ok
```

### Register/Deregister

```shell
micro register service '{"name": "foo", "version": "bar", "nodes": [{"id": "foo-1", "address": "127.0.0.1", "port": 8080}]}'
```

```shell
micro deregister service '{"name": "foo", "version": "bar", "nodes": [{"id": "foo-1", "address": "127.0.0.1", "port": 8080}]}'
```

## Proxy Remote Env

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

## Usage

```shell
NAME:
   micro - A cloud-native toolkit

USAGE:
   micro [global options] command [command options] [arguments...]
   
VERSION:
   0.8.0
   
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
    web		Run the micro web app

GLOBAL OPTIONS:
   --client 									Client for go-micro; rpc [$MICRO_CLIENT]
   --client_request_timeout 							Sets the client request timeout. e.g 500ms, 5s, 1m. Default: 5s [$MICRO_CLIENT_REQUEST_TIMEOUT]
   --client_retries "0"								Sets the client retries. Default: 1 [$MICRO_CLIENT_RETRIES]
   --client_pool_size "1"							Sets the client connection pool size. Default: 1 [$MICRO_CLIENT_POOL_SIZE]
   --client_pool_ttl 								Sets the client connection pool ttl. e.g 500ms, 5s, 1m. Default: 1m [$MICRO_CLIENT_POOL_TTL]
   --server_name 								Name of the server. go.micro.srv.example [$MICRO_SERVER_NAME]
   --server_version 								Version of the server. 1.1.0 [$MICRO_SERVER_VERSION]
   --server_id 									Id of the server. Auto-generated if not specified [$MICRO_SERVER_ID]
   --server_address 								Bind address for the server. 127.0.0.1:8080 [$MICRO_SERVER_ADDRESS]
   --server_advertise 								Used instead of the server_address when registering with discovery. 127.0.0.1:8080 [$MICRO_SERVER_ADVERTISE]
   --server_metadata [--server_metadata option --server_metadata option]	A list of key-value pairs defining metadata. version=1.0.0 [$MICRO_SERVER_METADATA]
   --broker 									Broker for pub/sub. http, nats, rabbitmq [$MICRO_BROKER]
   --broker_address 								Comma-separated list of broker addresses [$MICRO_BROKER_ADDRESS]
   --registry 									Registry for discovery. consul, mdns [$MICRO_REGISTRY]
   --registry_address 								Comma-separated list of registry addresses [$MICRO_REGISTRY_ADDRESS]
   --selector "cache"								Selector used to pick nodes for querying [$MICRO_SELECTOR]
   --server 									Server for go-micro; rpc [$MICRO_SERVER]
   --transport 									Transport mechanism used; http [$MICRO_TRANSPORT]
   --transport_address 								Comma-separated list of transport addresses [$MICRO_TRANSPORT_ADDRESS]
   --enable_acme								Enables ACME support via Let's Encrypt. ACME hosts should also be specified. [$MICRO_ENABLE_ACME]
   --acme_hosts 								Comma separated list of hostnames to manage ACME certs for [$MICRO_ACME_HOSTS]
   --enable_tls									Enable TLS support. Expects cert and key file to be specified [$MICRO_ENABLE_TLS]
   --tls_cert_file 								Path to the TLS Certificate file [$MICRO_TLS_CERT_FILE]
   --tls_key_file 								Path to the TLS Key file [$MICRO_TLS_KEY_FILE]
   --tls_client_ca_file 							Path to the TLS CA file to verify clients against [$MICRO_TLS_CLIENT_CA_FILE]
   --api_address 								Set the api address e.g 0.0.0.0:8080 [$MICRO_API_ADDRESS]
   --proxy_address 								Proxy requests via the HTTP address specified [$MICRO_PROXY_ADDRESS]
   --web_address 								Set the web UI address e.g 0.0.0.0:8082 [$MICRO_WEB_ADDRESS]
   --register_ttl "0"								Register TTL in seconds [$MICRO_REGISTER_TTL]
   --register_interval "0"							Register interval in seconds [$MICRO_REGISTER_INTERVAL]
   --api_handler 								Specify the request handler to be used for mapping HTTP requests to services; {api, proxy, rpc} [$MICRO_API_HANDLER]
   --api_namespace 								Set the namespace used by the API e.g. com.example.api [$MICRO_API_NAMESPACE]
   --web_namespace 								Set the namespace used by the Web proxy e.g. com.example.web [$MICRO_WEB_NAMESPACE]
   --api_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_API_CORS]
   --web_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_WEB_CORS]
   --proxy_cors 								Comma separated whitelist of allowed origins for CORS [$MICRO_PROXY_CORS]
   --enable_stats								Enable stats [$MICRO_ENABLE_STATS]
   --help, -h									show help
   --version									print the version
```
