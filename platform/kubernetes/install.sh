#!/bin/bash

# Run this script to install the platform on a kubernetes cluster. 

# NOTE: This script will not set the cloudflare or slack tokens in the secret. Hence, the 
# clients (web, api, proxy, bot) will have a status of "CreateContainerConfigError" until these
# secrets are manually added.

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
bash install.sh;
cd ../etcd;
bash install.sh;
cd ../nats;
bash install.sh;

# move to the /kubernetes folder and apply the deployments
cd ../..;
kubectl apply -f network

# go back to the top level
cd ..;