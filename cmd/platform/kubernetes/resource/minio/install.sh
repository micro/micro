#!/bin/bash

# move into the certs directory
cd certs;

# generate a certificate authority
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -;

# generate certificates for client and server
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client.json | cfssljson -bare client;
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server.json | cfssljson -bare server;

# create the secrets in minio
kubectl create secret generic minio-client-certs --from-file=ca.crt=ca.pem --from-file=public.crt=client.pem --from-file=private.key=client-key.pem;
kubectl create secret generic minio-server-certs --from-file=ca.pem --from-file=public.crt=server.pem --from-file=private.key=server-key.pem;

# move back into the minio directory
cd ../;

# add the nats helm chart
helm repo add minio https://helm.min.io/

helm install minio-cluster minio/minio \
  --set tls.enabled=true,tls.certSecret=minio-server-certs \
  --set persistence.size=3Gi \
  --set resources.requests.memory=512

# wait for the nats cluster to start before we start the streaming cluster
kubectl wait --for=condition=Ready pod/minio-cluster-0 --timeout=180s