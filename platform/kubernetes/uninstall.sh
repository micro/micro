#!/bin/bash

# Run this script to uninstall the platform from a kubernetes cluster

# reomve micro secrets
kubectl delete secret micro-secrets

# uninstall the resources
cd ./resource/cockroachdb;
bash uninstall.sh;
cd ../etcd;
bash uninstall.sh;
cd ../nats;
bash uninstall.sh;

# move to the /kubernetes folder and apply the deployments
cd ../..;
kubectl delete -f network

# go back to the top level
cd ..;