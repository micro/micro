---
title: Kubernetes Deployment
keywords: kubernetes
tags: [kubernetes]
sidebar: home_sidebar
permalink: /deploy-kubernetes
summary: 
---


This doc provides a guide to running micro on kubernetes.

## Getting Started

- [Dependencies](#dependencies)
- [Deployment](#deployment)
- [Micro API](#micro-api)
- [Micro Web](#micro-web)

## Dependencies

On kubernetes we would recommend running [etcd](https://github.com/etcd-io/etcd) and [nats](https://github.com/nats-io/nats-server).

- etcd is used for highly scalable service discovery
- NATS is used for asynchronous messaging

To install etcd ([instructions](https://github.com/helm/charts/tree/master/stable/etcd-operator))

```
helm install --name my-release --set customResources.createEtcdClusterCRD=true stable/etcd-operator
```

To install nats ([instructions](https://github.com/helm/charts/tree/master/stable/nats))

```
helm install my-release stable/nats
```

You should now have the required dependencies.

## Deployment

Here's an example k8s deployment for a micro service

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  namespace: default
  name: greeter
spec:
  replicas: 1
  selector:
    matchLabels:
      name: greeter-srv
      micro: service
  template:
    metadata:
      labels:
        name: greeter-srv
        micro: service
    spec:
      containers:
        - name: greeter
          command: [
		"/greeter-srv",
	  ]
          image: micro/go-micro
          imagePullPolicy: Always
          ports:
          - containerPort: 8080
            name: greeter-port
          env:
          - name: MICRO_SERVER_ADDRESS
            value: "0.0.0.0:8080"
          - name: MICRO_BROKER
            value: "nats"
          - name: MICRO_BROKER_ADDRESS
            value: "nats-cluster"
          - name: MICRO_REGISTRY
            value: "etcd"
          - name: MICRO_REGISTRY_ADDRESS
            value: "etcd-cluster-client"
```

Deploy with kubectl

```
kubectl apply -f greeter.yaml
```

## Micro API

To deploy the micro api use the following config. Note the ENABLE_ACME env var where you want Let's Encrypt SSL by default.

Create the api service

```
apiVersion: v1
kind: Service
metadata:
  name: micro-api
  namespace: default
  labels:
    name: micro-api
    micro: service
spec:
  ports:
  - name: https
    port: 443
    targetPort: 443
  selector:
    name: micro-api
    micro: service
  type: LoadBalancer
```

Create the deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: micro-api
  labels:
    micro: service
spec:
  replicas: 3
  selector:
    matchLabels:
      name: micro-api
      micro: service
  strategy:
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
  template:
    metadata:
      labels:
        name: micro-api
        micro: service
    spec:
      containers:
      - name: api
        env:
        - name: MICRO_ENABLE_STATS
          value: "true"
        - name: MICRO_BROKER
          value: "nats"
        - name: MICRO_BROKER_ADDRESS
          value: "nats-cluster"
        - name: MICRO_REGISTRY
          value: "etcd"
        - name: MICRO_REGISTRY_ADDRESS
          value: "etcd-cluster-client"
        - name: MICRO_REGISTER_TTL
          value: "60"
        - name: MICRO_REGISTER_INTERVAL
          value: "30"
        - name: MICRO_ENABLE_ACME
          value: "true"
        args:
        - api
        image: micro/micro
        imagePullPolicy: Always
        ports:
        - containerPort: 443
          name: api-port
```

## Micro Web

To deploy the micro web use the following config. Note the ENABLE_ACME env var where you want Let's Encrypt SSL by default.

Create the service

```
apiVersion: v1
kind: Service
metadata:
  name: micro-web
  namespace: default
  labels:
    name: micro-web
    micro: service
spec:
  ports:
  - name: https
    port: 443
    targetPort: 443
  selector:
    name: micro-web
    micro: service
  type: LoadBalancer
```

Create the deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: micro-web
  labels:
    micro: service
spec:
  replicas: 1
  selector:
    matchLabels:
      name: micro-web
      micro: service
  template:
    metadata:
      labels:
        name: micro-web
        micro: service
    spec:
      containers:
      - name: web
        env:
        - name: MICRO_BROKER
          value: "nats"
        - name: MICRO_BROKER_ADDRESS
          value: "nats-cluster"
        - name: MICRO_ENABLE_STATS
          value: "true"
        - name: MICRO_REGISTRY
          value: "etcd"
        - name: MICRO_REGISTRY_ADDRESS
          value: "etcd-cluster-client"
        - name: MICRO_ENABLE_ACME
          value: "true"
        args:
        - web
        image: micro/micro
        imagePullPolicy: Always
        ports:
        - containerPort: 443
          name: web-port
```

