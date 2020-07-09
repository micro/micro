#!/bin/bash

# uninstall the cluster 
kubectl delete -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml;
kubectl delete -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml;
kubectl delete -f nats.yaml;

# delete the secrets
kubectl delete secret nats-client-certs;
kubectl delete secret nats-server-certs;
kubectl delete secret nats-peer-certs;
