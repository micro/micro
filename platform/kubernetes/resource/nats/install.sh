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

# install the cluster 
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml;
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml;
kubectl wait --timeout=180s -n default --for=condition=available deployment/nats-operator
kubectl apply -f nats.yaml;