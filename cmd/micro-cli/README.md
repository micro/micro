# Micro CLI

An interactive cli prompt for Micro

## Usage

Install it

```
go get github.com/micro/micro/v5/cmd/micro-cli@master
```

Run it

```
micro-cli
```

Should return the prompt

```
micro>
```

Type `help` for commands. Supports same commands as the normal CLI.

# Micro CLI Admin Commands

The Micro CLI provides an interactive shell and direct commands for managing go-micro platform primitives and services. All commands use consistent naming, input/output, and error handling.

## Features

- **Store**: (see interactive shell or service calls)
- **Broker**:
  - `broker publish [topic] [message]` — Publish a message
  - `broker subscribe [topic]` — Subscribe to a topic (one message)
- **Config**:
  - `config get [key]` — Get a config value
  - `config set [key] [value]` — Set a config value
  - `config delete [key]` — Delete a config value
  - `config list` — List all config keys
- **Registry**:
  - `registry list` — List all services
  - `registry get [name]` — Get a service
  - `registry register [name] [node_id] [address] [version]` — Register a service
  - `registry deregister [name] [node_id]` — Deregister a service
- **Service Calls**:
  - `call [service] [endpoint] [request]` — Call any service endpoint
  - `describe [service]` — Describe a service and its endpoints
  - `services` — List all services

## Usage

Start the CLI:

```sh
micro cli
```

You will see a prompt:

```
micro> 
```

Type any command, e.g.:

```
micro> broker publish demo hello
micro> config set foo bar
micro> registry list
micro> call helloworld Helloworld.Call '{"name":"Alice"}'
```

Type `help` or `?` for a list of commands.

## Error Handling

All commands print clear error messages on failure.

## See Also
- [Web Admin UI](../micro-web/README.md)
- [API Admin Endpoints](../micro-api/README.md)

---

For more information, see the main [Micro documentation](../../README.md).
