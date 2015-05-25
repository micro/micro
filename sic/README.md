# Micro SideCar

The sidecar provides features to integrate any application into the Micro ecosystem.

It is similar to Netflix's sidecar called [Prana](https://github.com/Netflix/Prana)

## Features

- Registration with discovery system
- Host discovery of other services
- Health checking of services
- HTTP API and load balancing requests
- PubSub via WebSockets

## Getting Started

### Install

```shell
$ go get github.com/myodc/micro
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
{
	"name":"go.micro.service.template",
	"nodes":[{
		"id":"go.micro.service.template-c5718d29-da2a-11e4-be11-68a86d0d36b6",
		"address":"[::]","port":60728
	}]
}
```

### Healthchecking

Start micro sidecar with "--healthcheck_url=" to enable the healthchecker

```shell
$ micro sidecar --server_name=foo --server_address=127.0.0.1:9090 \
	--healthcheck_url=http://127.0.0.1:9090/_status/health
I0523 12:25:36.229536   85658 sic.go:184] Registering foo-6ebf29c0-013e-11e5-b55f-68a86d0d36b6
I0523 12:25:36.241680   85658 sic.go:188] Starting sidecar healthchecker
```

### HTTP RPC API

Query micro services via the http rpc api.

```shell
$ curl  -d 'service=go.micro.service.template' \
	-d 'method=Example.Call' \
	-d 'request={"name": "John"}' http://127.0.0.1:8081/rpc
{"msg":"go.micro.service.template-c5718d29-da2a-11e4-be11-68a86d0d36b6: Hello John"}
```

### PubSub via WebSockets

Connect to the micro pub/sub broker via a websocket interface

```go
c, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:8081/broker?topic=foo", make(http.Header))

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
