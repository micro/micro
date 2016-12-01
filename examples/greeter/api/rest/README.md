# Go Restful API Example

This is an example of how to serve REST behind the API using go-restful

## Getting Started

### Run the Micro API

```
$ micro api --handler=proxy
```

### Run the Greeter Service

```
$ go run greeter/server/main.go
```

### Run the Greeter API

```
$ go run go-restful.go
Listening on [::]:64738
```

### Curl the API

Test the index
```
curl http://localhost:8080/greeter
{
  "message": "Hi, this is the Greeter API"
}
```

Test a resource
```
 curl http://localhost:8080/greeter/asim
{
  "msg": "Hello asim"
}
```
