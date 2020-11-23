#!/bin/bash

# Run this script to uninstall the platform from a kubernetes cluster

# reomve micro secrets
kubectl delete secret micro-secrets

# uninstall the resources
for d in ./resource/*/; do
  pushd $d
  MICRO_ENV=$ENV bash uninstall.sh
  popd
done

# delete the deployments and services
kubectl delete -f ./service

# go back to the top level
cd ..;
