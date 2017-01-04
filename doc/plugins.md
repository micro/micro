# Plugins

Micro is a pluggable toolkit and framework. The internal features can be swapped out with any [go-plugins](https://github.com/micro/go-plugins).

## Usage

Plugins can be added to go-micro in the following ways. By doing so they'll be available to set via command line args or environment variables.

### Import Plugins

```go
import (
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)

func main() {
	// Parse CLI flags
	cmd.Init()
}
```

The same is achieved when calling ```service.Init```

```go
import (
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)

func main() {
	service := micro.NewService(
		// Set service name
		micro.Name("my.service"),
	)

	// Parse CLI flags
	service.Init()
}
```

### Use via CLI Flags

Activate via a command line flag

```shell
go run service.go --broker=rabbitmq --registry=kubernetes --transport=nats
```

### Use Plugins Directly

CLI Flags provide a simple way to initialise plugins but you can do the same yourself.

```go
import (
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	registry := kubernetes.NewRegistry() //a default to using env vars for master API

	service := micro.NewService(
		// Set service name
		micro.Name("my.service"),
		// Set service registry
		micro.Registry(registry),
	)
}
```

## Build Pattern

You may want to swap out plugins using automation or add plugins to the micro toolkit. 
An easy way to do this is by maintaining a separate file for plugin imports and including it during the build.

Create file plugins.go
```go
package main

import (
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)
```

Build with plugins.go
```shell
go build -o service main.go plugins.go
```

Run with plugins
```shell
service --broker=rabbitmq --registry=kubernetes --transport=nats
```

## Rebuild Toolkit With Plugins

If you want to integrate plugins simply link them in a separate file and rebuild

Create a plugins.go file
```go
import (
        // etcd v3 registry
        _ "github.com/micro/go-plugins/registry/etcdv3"
        // nats transport
        _ "github.com/micro/go-plugins/transport/nats"
        // kafka broker
        _ "github.com/micro/go-plugins/broker/kafka"
```

Build binary
```shell
// For local use
go build -i -o micro ./main.go ./plugins.go
// For docker image
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o micro ./main.go ./plugins.go
```

Flag usage of plugins
```shell
micro --registry=etcdv3 --transport=nats --broker=kafka
```
