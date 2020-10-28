# Tunnel

The micro tunnel is a tunneling interface used to interconnect multiple networks. 

## Usage

The micro tunnel can be used as a go-micro package or as the `micro tunnel` service

```
t := tunnel.NewTunnel(
	tunnel.Nodes(...) // list of nodes to connect to
)

// connect the tunnel
err := t.Connect()
// close the tunnel
defer t.Close()

// listen for requests
l, err := t.Listen(addr)

for {
    // accept messages
    c, err := l.Accept()
    // do something
}

// dial endpoint
c, err := t.Dial(addr)
```

As a service

```
# Run the tunnel server
micro tunnel

# connect to the tunnel specifying the server
micro tunnel --address=:8090 --server=:8083

# Use the tunnel as a proxy
MICRO_PROXY=go.micro.tunnel go run myservice.go
```

## Protocol

The protocol is state machine as defined by the Micro-Tunnel header. 

Possible message types

```
connect
close
keepalive
open
listen
message
```
