#!/bin/bash

# which env to deploy to. does not yet switch k8s context
ENV=$1
# default to size for platform
DB_SIZE=100Gi

# Run this script to install the platform on a kubernetes cluster. 

# NOTE: This script will not set the cloudflare or slack tokens in the secret. Hence, the 
# clients (web, api, proxy, bot) will have a status of "CreateContainerConfigError" until these
# secrets are manually added.

# expect an env to be specified
if [ "$ENV" == "" ]; then
  echo "Must specify env e.g ./install.sh {dev|staging|platform}"
  exit 1
fi

## Set DB to smaller size for staging
if [ "$ENV" != "platform" ]; then
  DB_SIZE=25Gi
fi

# Generate keys for JWT auth.
ssh-keygen -f /tmp/jwt -m pkcs8 -q -N "";
ssh-keygen -f /tmp/jwt -e  -m pkcs8 > /tmp/jwt.pub;
cat /tmp/jwt | base64 > /tmp/jwt-base64
cat /tmp/jwt.pub | base64 > /tmp/jwt-base64.pub

# Create the k8s secret
kubectl create secret generic micro-secrets \
  --from-file=auth_public_key=/tmp/jwt-base64.pub \
  --from-file=auth_private_key=/tmp/jwt-base64;

# Remove the files from tmp
rm /tmp/jwt /tmp/jwt.pub /tmp/jwt-base64 /tmp/jwt-base64.pub

# install the resources
cd ./resource/cockroachdb;
bash install.sh $DB_SIZE;
cd ../etcd;
bash install.sh;
cd ../nats;
bash install.sh;

# move to the /kubernetes folder and apply the deployments
cd ../..;

# replace m3o.com with m3o.dev
if [ $ENV == "staging" ]; then
  sed -i 's@m3o.com@m3o.dev@g' network/*.yaml
fi

# execute the yaml
kubectl apply -f network

# replace back
if [ $ENV == "staging" ]; then
  sed -i 's@m3o.dev@m3o.com@g' network/*.yaml
fi

# go back to the top level
cd ..;
