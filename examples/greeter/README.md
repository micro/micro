# Greeter Service

An example Go service running with go-micro

## What's here?

- **Server** - an RPC greeter service
- **Client** - an RPC client that calls the service once
- **Api** - examples of RPC API and RESTful API
- **Web** - how to use go-web to write web services

### Prerequisites

Service discovery is required either running Consul as the default or MDNS for zero dependencies. You can also use plugins from [micro/plugins](https://github.com/micro/go-plugins).

**MDNS**

To use MDNS just add the flag `--registry=mdns` on the command line.

**Consul**

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run Service

Start go.micro.srv.greeter
```shell
$ go run server/main.go
I0523 12:27:56.685569   90096 server.go:113] Starting server go.micro.srv.greeter id go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6
I0523 12:27:56.685890   90096 rpc_server.go:112] Listening on [::]:52019
I0523 12:27:56.685931   90096 server.go:95] Registering node: go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6
```

### Client

Call go.micro.srv.greeter via client
```shell
$ go run client/client.go
go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6: Hello John
```

Examples of client usage via other languages can be found in the client directory.

### API

HTTP based requests can be made via the micro API. Micro logically separates API services from backend services. By default the micro API 
accepts HTTP requests and converts to *api.Request and *api.Response types. Find them here [micro/api/proto](https://github.com/micro/micro/tree/master/api/proto).

Run the go.micro.api.greeter API Service
```shell
$ go run api/api.go 
I0523 12:27:25.475548   89125 server.go:113] Starting server go.micro.api.greeter id go.micro.api.greeter-afdcc369-013e-11e5-8710-68a86d0d36b6
I0523 12:27:25.475890   89125 rpc_server.go:112] Listening on [::]:51997
I0523 12:27:25.475927   89125 server.go:95] Registering node: go.micro.api.greeter-afdcc369-013e-11e5-8710-68a86d0d36b6
```

Run the micro API
```shell
$ micro api
I0523 12:23:23.413940   81384 api.go:131] API Rpc handler /rpc
I0523 12:23:23.414238   81384 api.go:143] Listening on [::]:8080
I0523 12:23:23.414272   81384 server.go:113] Starting server go.micro.api id go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
I0523 12:23:23.414355   81384 rpc_server.go:112] Listening on [::]:51938
I0523 12:23:23.414399   81384 server.go:95] Registering node: go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
```

Call go.micro.api.greeter via API
```shell
curl http://localhost:8080/greeter/say/hello?name=John
{"message":"Hello John"}
```

Examples of other API handlers can be found in the API directory.

### Sidecar

The sidecar is a language agnostic RPC proxy.

Run the micro sidecar
```shell
$ micro sidecar
2016/12/01 14:26:44 Registering Root Handler at /
2016/12/01 14:26:44 Registering Registry handler at /registry
2016/12/01 14:26:44 Registering RPC handler at /rpc
2016/12/01 14:26:44 Registering Broker handler at /broker
2016/12/01 14:26:44 Listening on [::]:8081
2016/12/01 14:26:44 Listening on [::]:51287
2016/12/01 14:26:44 Broker Listening on [::]:51288
2016/12/01 14:26:44 Registering node: go.micro.sidecar-2f727f0c-b7d2-11e6-aa3b-68a86d0d36b6
```

Call go.micro.srv.greeter via sidecar
```shell
curl -H 'Content-Type: application/json' -d '{"name": "john"}' http://localhost:8081/greeter/say/hello
{"msg":"Hello john"}
```

The sidecar provides all the features of go-micro as a HTTP API. Learn more at [micro/car](https://github.com/micro/micro/tree/master/car).

### Web Usage

The micro web is a web dashboard and reverse proxy to run web apps as microservices.

Run go.micro.web.greeter
```
$ go run web/web.go 
Listening on [::]:51348
```

Run the micro web
```shell
$ micro web
2016/12/01 14:30:22 Listening on [::]:8082
2016/12/01 14:30:22 Listening on [::]:51334
2016/12/01 14:30:22 Broker Listening on [::]:51335
2016/12/01 14:30:22 Registering node: go.micro.web-b10c7e52-b7d2-11e6-a39c-68a86d0d36b6
```

Browse to http://localhost:8082/greeter
