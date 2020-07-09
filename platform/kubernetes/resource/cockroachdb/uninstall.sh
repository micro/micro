#!/bin/bash

# delete the cluster using helm
helm delete cockroachdb-cluster;

# delete the secrets 
kubectl delete secret cockroachdb-client-certs;
kubectl delete secret cockroachdb-server-certs;
kubectl delete secret cockroachdb-peer-certs;
