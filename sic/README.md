# Micro SideCar

The sidecar provides features to integrate any application into the Micro ecosystem.

It is similar to Netflix's sidecar called [Prana](https://github.com/Netflix/Prana)

## Features

- Registration with discovery system
- Host discovery of other services
- Health checking of services
- HTTP API and load balancing requests
- Access to basic key/value store
- PubSub via WebSockets

## Getting Started

### Install

```shell
$ go get github.com/asim/micro
```

### Run

The micro sidecar runs on port 8081 by default.

Start the sidecar. Specify your app server name and address.

```shell
$ micro sidecar --server_name=foo --server_address=127.0.0.1:9090
```

### Host Discovery

```shell
curl http://127.0.0.1:8081/registry?service=go.micro.service.template
{"name":"go.micro.service.template","nodes":[{"id":"go.micro.service.template-c5718d29-da2a-11e4-be11-68a86d0d36b6","address":"[::]","port":60728}]}
```

### Healthchecking

Start micro sidecar with "--healthcheck_url=" to enable the healthchecker

```shell
$ micro sidecar --server_name=foo --server_address=127.0.0.1:9090 --healthcheck_url=http://127.0.0.1:9090/_status/health
I0409 20:45:53.430312   27577 sic.go:182] Registering foo-08378009-def1-11e4-a015-68a86d0d36b6
***I0409 20:45:53.437452   27577 sic.go:186] Starting sidecar healthchecker***
```

### Query basic key/value store

Put item
```shell
$ curl -d 'key=foo' -d 'value=bar' http://127.0.0.1:8081/store
```

Get item
```shell
$ curl http://127.0.0.1:8081/store?key=foo
bar
```

Del
```shell
$ curl -XDELETE http://127.0.0.1:8081/store?key=foo
```

### HTTP RPC API

Query micro services via the http rpc api.

```shell
$ curl -d 'service=go.micro.service.template' -d 'method=Example.Call' -d 'request={"name": "John"}' http://127.0.0.1:8081/rpc
{"msg":"go.micro.service.template-c5718d29-da2a-11e4-be11-68a86d0d36b6: Hello John"}
```

### PubSub via WebSockets

Connect to the micro pub/sub broker via a websocket interface

```go
c, _, _ := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:8081/broker?topic=foo", mqServer, topic), make(http.Header))

go func() {
	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			return
		}
		var msg *broker.Message
		json.Unmarshal(p, &msg)
		fmt.Println(msg.Data)
	}
}()

ticker := time.NewTicker(time.Second)

for _ = range ticker.C {
	if err := c.WriteMessage(1, []byte(`hello world`)); err != nil {
		return
	}
}
```
