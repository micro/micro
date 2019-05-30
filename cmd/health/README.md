# Health

Health is a healthchecking sidecar

It exposes `/health` as a http endpoint and calls the rpc endpoint `Debug.Health`

Every go-micro service exposes `Debug.Health`. We can use this for liveness checks.

## Usage

### Install

```
go get github.com/micro/util/cmd/health
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

### Kubernetes

We can add the healthchecking sidecar to a service with a livenessProbe

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  namespace: default
  name: greeter
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: greeter-srv
        micro: go.micro.srv.greeter
    spec:
      containers:
        - name: greeter
          command: [
		"/greeter-srv",
		"--server_address=0.0.0.0:9091",
		"--broker_address=0.0.0.0:10001"
	  ]
          image: microhq/greeter-srv:kubernetes
          imagePullPolicy: Always
          ports:
          - containerPort: 9091
            name: greeter-port
        - name: liveness
          command: [
		"/health",
		"--server_name=greeter",
		"--server_address=localhost:9091"
	  ]
          image: microhq/health:kubernetes
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 3
            periodSeconds: 3
```
