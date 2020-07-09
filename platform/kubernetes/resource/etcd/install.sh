#!/bin/bash

# move into the certs directory
cd certs;

# generate a certificate authority
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -

# generate certificates for client, server and peer
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client.json | cfssljson -bare client;
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server.json | cfssljson -bare server;
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=peer peer.json | cfssljson -bare peer;

# create the secrets in etcd
kubectl create secret generic etcd-client-certs --from-file=ca.crt=ca.pem --from-file=cert.pem=client.pem --from-file=key.pem=client-key.pem;
kubectl create secret generic etcd-server-certs --from-file=ca.crt=ca.pem --from-file=cert.pem=server.pem --from-file=key.pem=server-key.pem;
kubectl create secret generic etcd-peer-certs --from-file=ca.crt=ca.pem --from-file=cert.pem=peer.pem --from-file=key.pem=peer-key.pem;

# move back into the /etcd directory
cd ../;

# install the cluster using helm
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install etcd-cluster bitnami/etcd \
	--set auth.rbac.enabled=false \
	--set auth.peer.secureTransport=true \
	--set auth.peer.enableAuthentication=true \
	--set auth.peer.existingSecret=etcd-peer-certs \
	--set auth.client.secureTransport=true \
	--set auth.client.enableAuthentication=true \
	--set auth.client.existingSecret=etcd-server-certs \
	--set statefulset.replicaCount=1 \
	--set livenessProbe.enabled=false \
	--set readinessProbe.enabled=false