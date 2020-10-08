#!/bin/bash
# REQUIRED MICRO ENV S3_ACCESS_KEY S3_SECRET_KEY S3_ENDPOINT S3_REGION

# generate a set of credentials
accessKey=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 20 ; echo)
secretKey=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 40 ; echo)
region=$S3_REGION
if [[ ! $region ]]; then
  region=fr-par # default
fi
# create the secret to store the creds
kubectl create secret generic minio-creds --from-literal=access_key=$accessKey --from-literal=secret_key=$secretKey --from-literal=region=$region

# add the nats helm chart
helm repo add minio https://helm.min.io/

if [[ $MICRO_ENV == "dev" ]]; then
  overrides="--set persistence.size=3Gi --set s3gateway.enabled=false"
else
  overrides="--set s3gateway.serviceEndpoint=$S3_ENDPOINT --set s3gateway.accessKey=$S3_ACCESS_KEY --set s3gateway.secretKey=$S3_SECRET_KEY"
fi

helm install minio-cluster minio/minio --version 7.1.2 $overrides -f values.yaml \
  --set accessKey=$accessKey \
  --set secretKey=$secretKey \
  --set resources.requests.memory=512

# wait for the minio cluster to start
kubectl wait deployment minio-cluster --timeout=60s -n default --for=condition=available