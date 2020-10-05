#!/bin/bash

# uninstall the cluster using helm
helm delete minio-cluster

# delete the secrets 
kubectl delete secret minio-creds;