# Nats

Nats is used for message

## Usage

First run the following

```
## Install NATS Operator
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml
```

Then apply the k8s CRD

```
kubectl apply -f nats.yaml
```
