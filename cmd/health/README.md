# Health

Health is a healthchecking sidecar

It exposes `/health` as a http endpoint and calls the rpc endpoint `Debug.Health`

Every go-micro service exposes `Debug.Health`. We can use this for liveness checks.

## Usage

### Install

```
go get github.com/micro/micro/cmd/health
```

or

```
docker pull microhq/health
```

### Run

Run the healthchecker specifying the service name and address to check

```
health --server_name=greeter --server_address=localhost:9091
```

### Call

Call the healthchecker on localhost:8080

```
curl http://localhost:8080/health
```

Response will be 200 OK

```
OK
```

