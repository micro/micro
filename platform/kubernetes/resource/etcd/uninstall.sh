#!/bin/bash

# uninstall the cluster using helm
helm delete etcd-cluster

# delete the secrets 
kubectl delete secret etcd-client-certs;
kubectl delete secret etcd-server-certs;
kubectl delete secret etcd-peer-certs;