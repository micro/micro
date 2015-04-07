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
