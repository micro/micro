# Kubernetes Deployment

This repo serves as the kubernetes deployment for the platform.

## Overview

The platform consists of the following

- **infra** - shared infrastructure dependencies that must run in every region
- **control** - the control plane run only in DO
- **runtime** - the micro runtime run ontop of the shared infra
- **services** - adhoc services we run on the platform

## Dependencies

We have dependencies to get started

- Kubectl
- Kubectx
- Helm
- Etcd
- NATS

## Usage

Check the [GCP](gcloud.md) or [AWS](aws.md) docs for specific instructions. Or:

1. Spin up managed K8s somewhere

2. Setup local env
  - Install `kubectl` https://kubernetes.io/docs/tasks/tools/install-kubectl/
  - Install `kubectx` https://github.com/ahmetb/kubectx
  - Install `helm` https://github.com/helm/helm
    * See [this](https://helm.sh/docs/using_helm/#tiller-namespaces-and-rbac)
    * `kubectl -n kube-system create serviceaccount tiller`
    * `kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller`
    * `helm init --service-account=tiller`

3. Install `etcd`
  - `helm repo update`
  - `helm install stable/etcd-operator --version="0.10.0" --set customResources.createEtcdClusterCRD=true --set etcdCluster.version="3.3.15" --set etcdCluster.image.tag="v3.3.15"`
  - read the docs [here](https://etcd.io/docs/v3.3.12/)

4. Install `nats`
 - `kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml`
 - `kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml`
 - `kubectl apply -f network/config/kubernetes/services/infra/nats.yaml`

4. Install Micro core
  - kubectl apply -f ../kubernetes
  - Create external load balancers https://www.digitalocean.com/docs/kubernetes/how-to/add-load-balancers/

5. Install `etcd` in your local environment and query remote `etcd` cluster
  - `kubectl port-forward service/etcd-cluster-client -n default 2379`
  - `ETCDCTL_API=3 etcdctl version`
  - `ETCDCTL_API=3 etcdctl -w table member list`

6. Accessing particular `etcd` node directly
  - `kubectl exec -it ETCD_NODE_POD_NAME -- /bin/sh`

## Ports

- API - 443 (UDP)
- Web - 443 (UDP)
- Network - 8085 (UDP)
- Tunnel - 8083 (UDP)

## DNS Records

- micro.mu -> seed node
- api.micro.mu -> micro api
- web.micro.mu -> micro web
- tunnel.micro.mu -> micro tunnel
- network.micro.mu -> micro network

## Healing Etcd

- `helm upgrade [release-name] stable/etcd-operator --set etcdCluster.size=3 --set etcdCluster.image.tag="v3.4.2" --set customResources.createEtcdClusterCRD=truestable/etcd-operator`
