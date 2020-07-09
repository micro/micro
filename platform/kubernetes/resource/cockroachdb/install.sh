#!/bin/bash

# move into the certs directory
cd certs;

# generate a certificate authority
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -

# generate certificates for client, server and peer
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client.json | cfssljson -bare client;
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server.json | cfssljson -bare server;
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=peer peer.json | cfssljson -bare peer;

# create the secrets in cockroachdb
kubectl create secret generic cockroachdb-client-certs --from-file=ca.crt=ca.pem --from-file=cert.pem=client.pem --from-file=key.pem=client-key.pem;
kubectl create secret generic cockroachdb-server-certs --from-file=ca.crt=ca.pem --from-file=tls.crt=server.pem --from-file=tls.key=server-key.pem;
kubectl create secret generic cockroachdb-peer-certs --from-file=ca.crt=ca.pem --from-file=tls.crt=peer.pem --from-file=tls.key=peer-key.pem;

# move back into the /cockroachdb directory
cd ../;

# install the cluster using helm
helm repo add cockroachdb https://charts.cockroachdb.com/
helm install cockroachdb-cluster cockroachdb/cockroachdb \
  --set statefulset.replicas=1 \
  --set storage.persistentVolume.size=10gi \
  --set tls.certs.clientRootSecret=cockroachdb-peer-certs \
  --set tls.certs.nodeSecret=cockroachdb-server-certs \
  --set tls.certs.tlsSecret=true \
  --set tls.certs.provided=true \
  --set tls.enabled=true