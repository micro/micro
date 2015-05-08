# Greeter Service

An example Go service running with go-micro

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
```

Run Service
```
$ go run server/main.go

I0407 20:47:39.267559   21462 rpc_server.go:156] Rpc handler /_rpc
I0407 20:47:39.276808   21462 server.go:90] Starting server go.micro.srv.greeter id go.micro.srv.greeter-f2797eaf-dd5e-11e4-b169-68a86d0d36b6
I0407 20:47:39.276914   21462 rpc_server.go:187] Listening on [::]:64343
I0407 20:47:39.276986   21462 server.go:76] Registering go.micro.srv.greeter-f2797eaf-dd5e-11e4-b169-68a86d0d36b6
```

Test Service
```
$ go run client/client.go

go.micro.srv.greeter-165e25c3-dd5f-11e4-a8ed-68a86d0d36b6: Hello John
```

### REST Micro API usage

You can also construct REST based requests via an API service and the Micro API. REST based API service handlers take 
*api.Request and *api.Response types which can be found in [github.com/myodc/micro/api/proto](https://github.com/myodc/micro/tree/master/api/proto)

Run the API Service
```
$ go run api/api.go 
I0509 00:24:15.163051   23534 rpc_server.go:156] Rpc handler /_rpc
I0509 00:24:15.165418   23534 server.go:90] Starting server go.micro.api.greeter id go.micro.api.greeter-577001b7-f5d9-11e4-aaa8-68a86d0d36b6
I0509 00:24:15.165504   23534 rpc_server.go:187] Listening on [::]:50116
I0509 00:24:15.165535   23534 server.go:76] Registering go.micro.api.greeter-577001b7-f5d9-11e4-aaa8-68a86d0d36b6
```

Run the Micro API
```
$ micro api
I0509 00:25:10.512953   23577 rpc_server.go:156] Rpc handler /_rpc
I0509 00:25:10.515765   23577 api.go:129] API Rpc handler /rpc
I0509 00:25:10.515850   23577 api.go:141] Listening on [::]:8080
I0509 00:25:10.515870   23577 server.go:90] Starting server go.micro.api id go.micro.api-786dc223-f5d9-11e4-9c8b-68a86d0d36b6
I0509 00:25:10.515966   23577 rpc_server.go:187] Listening on [::]:50125
I0509 00:25:10.516006   23577 server.go:76] Registering go.micro.api-786dc223-f5d9-11e4-9c8b-68a86d0d36b6
```

Test via Curl
```
curl http://localhost:8080/greeter/say/hello?name=John
{"message":"go.micro.srv.greeter-8dfa9b29-f5d9-11e4-ba70-68a86d0d36b6: Hello John"}
```
