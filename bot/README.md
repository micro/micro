# Micro Bot

The micro bot is a bot that sits inside your microservices platform which you can interact with via Slack, HipChat, XMPP, etc. 
It mimics the functions of the CLI via messaging.

## Supported Inputs

- Slack

## Usage

### Run

```shell
$ micro bot --inputs=slack --slack_token=SLACK_TOKEN
```

### Help

In slack
```shell
$ micro help
ping - Returns pong
list services - Returns a list of registered services
get service [name] - Returns a registered service
health [service] - Returns health of a service
query [service] [method] [request] - Returns the response for a service query
register service [definition] - Registers a service
deregister service [definition] - Deregisters a service
hello - Returns a greeting
```
