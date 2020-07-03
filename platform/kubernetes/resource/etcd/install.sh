#!/bin/bash


#helm repo add bitnami https://charts.bitnami.com/bitnami
#helm install etcd-cluster bitnami/etcd \
#	--set auth.rbac.enabled=false \
#	--set auth.client.enableAuthentication=false \
#	--set livenessProbe.enabled=false \
#	--set readinessProbe.enabled=false

helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm install etcd-cluster stable/etcd-operator \
	--set resources.limits.cpu=500m \
	--set resources.limits.memory=512Mi \
	--set cluster.version=v3.4.9 \
	--set cluster.size=3 \
	--set customResources.createEtcdClusterCRD=true
