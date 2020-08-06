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
which ssh-keygen > /dev/null
if [ $? -eq 1 ]; then
  echo "Missing ssh-keygen command"
  exit 1
fi

which openssl > /dev/null
if [ $? -eq 1 ]; then
  echo "Missing openssl command"
fi

# generate new PEM key
ssh-keygen -t rsa -b 2048 -m PEM -f /tmp/jwt -q -N "";
# Don't add passphrase
openssl rsa -in /tmp/jwt -pubout -outform PEM -out /tmp/jwt.pub
# Base64 encode
base64 /tmp/jwt > /tmp/jwt-base64
base64 /tmp/jwt.pub > /tmp/jwt-base64.pub

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
  sed -i 's@m3o.com@m3o.dev@g' service/*.yaml
fi

# execute the yaml
kubectl apply -f service

# replace back
if [ $ENV == "staging" ]; then
  sed -i 's@m3o.dev@m3o.com@g' service/*.yaml
fi

# go back to the top level
cd ..;
