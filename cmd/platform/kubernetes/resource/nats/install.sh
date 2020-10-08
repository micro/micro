#!/bin/bash

# move into the certs directory
cd certs;

# generate a certificate authority
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -;

# generate certificates for client and server
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client.json | cfssljson -bare client;
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server.json | cfssljson -bare server;
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=peer peer.json | cfssljson -bare peer;

# create the secrets in nats
kubectl create secret generic nats-client-certs --from-file=ca.crt=ca.pem --from-file=cert.pem=client.pem --from-file=key.pem=client-key.pem;
kubectl create secret generic nats-server-certs --from-file=ca.pem --from-file=server.pem --from-file=server-key.pem;
kubectl create secret generic nats-peer-certs --from-file=ca.pem --from-file=route.pem=peer.pem --from-file=route-key.pem=peer-key.pem;

# move back into the nats directory
cd ../;

if [[ $MICRO_ENV == "dev" ]]; then
  overrides="--set cluster.enabled=false --set stan.replicas=1 --set nats.cluster.enabled=false --set store.cluster.enabled=false"
fi

# add the nats helm chart
helm repo add nats https://nats-io.github.io/k8s/helm/charts/

helm install nats-cluster nats/nats --version 0.5.6 -f nats-values.yaml $overrides 

# wait for the nats cluster to start before we start the streaming cluster
kubectl wait --for=condition=Ready pod/nats-cluster-0 --timeout=180s

helm install nats-streaming-cluster nats/stan --version 0.5.6 -f stan-values.yaml $overrides

# wait for the nats streaming cluster to start before we exit
kubectl wait --for=condition=Ready pod/nats-streaming-cluster-0 --timeout=180s 
