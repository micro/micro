#!/bin/bash

# uninstall the cluster 
helm uninstall nats-streaming-cluster;
helm uninstall nats-cluster;

# delete the secrets
kubectl delete secret nats-client-certs;
kubectl delete secret nats-server-certs;
kubectl delete secret nats-peer-certs;
