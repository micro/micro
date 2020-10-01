#!/bin/bash

# generate a set of credentials
accessKey=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 20 ; echo)
secretKey=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 40 ; echo)

# create the secret to store the creds
kubectl create secret generic minio-creds --from-literal=access_key=$accessKey --from-literal=secret_key=$secretKey

# add the nats helm chart
helm repo add minio https://helm.min.io/

helm install minio-cluster minio/minio \
  --set accessKey=$accessKey \
  --set secretKey=$secretKey \
  --set tls.certSecret=minio-server-certs \
  --set persistence.size=3Gi \
  --set resources.requests.memory=512

# wait for the minio cluster to start
kubectl wait deployment minio-cluster --timeout=60s -n default --for=condition=available