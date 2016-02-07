# Greeter Service

An example Go service running with go-micro

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

Run Service
```
$ go run server/main.go
I0523 12:27:56.685569   90096 server.go:113] Starting server go.micro.srv.greeter id go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6
I0523 12:27:56.685890   90096 rpc_server.go:112] Listening on [::]:52019
I0523 12:27:56.685931   90096 server.go:95] Registering node: go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6
```

Test Service
```
$ go run client/client.go
go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6: Hello John
```

### REST Micro API usage

You can also construct REST based requests via an API service and the Micro API. REST based API service handlers take 
*api.Request and *api.Response types which can be found in [github.com/micro/micro/api/proto](https://github.com/micro/micro/tree/master/api/proto)

Run the API Service
```
$ go run api/api.go 
I0523 12:27:25.475548   89125 server.go:113] Starting server go.micro.api.greeter id go.micro.api.greeter-afdcc369-013e-11e5-8710-68a86d0d36b6
I0523 12:27:25.475890   89125 rpc_server.go:112] Listening on [::]:51997
I0523 12:27:25.475927   89125 server.go:95] Registering node: go.micro.api.greeter-afdcc369-013e-11e5-8710-68a86d0d36b6
```

Run the Micro API
```
$ micro api
I0523 12:23:23.413940   81384 api.go:131] API Rpc handler /rpc
I0523 12:23:23.414238   81384 api.go:143] Listening on [::]:8080
I0523 12:23:23.414272   81384 server.go:113] Starting server go.micro.api id go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
I0523 12:23:23.414355   81384 rpc_server.go:112] Listening on [::]:51938
I0523 12:23:23.414399   81384 server.go:95] Registering node: go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
```

Test API via Curl
```
curl http://localhost:8080/greeter/say/hello?name=John
{"message":"go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6: Hello John"}
```

Test RPC via Curl
```
curl -d 'service=go.micro.srv.greeter' -d 'method=Say.Hello' -d 'request={"name": "john"}' http://localhost:8080/rpc
{"msg":"go.micro.srv.greeter-c2770c77-013e-11e5-b4d6-68a86d0d36b6: Hello john"}
```

