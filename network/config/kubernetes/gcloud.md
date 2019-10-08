# Setting up Micro on Google Cloud Platform

## Prerequisites 

Install Google cloud SDK using your favourite method, e.g.
```shell
brew cask install google-cloud-sdk
```

Log in to Google Cloud and select your project, region, etc.
```shell
gcloud init # --console-only if in a remote shell
```

Install `kubectl`, e.g.
```shell
brew install kubernetes-cli
```

N.B. `gcloud` has an interactive shell if you wish to use it:

```shell
gcloud components install beta
gcloud beta interactive
```

## Create a GKE Cluster

First, set up a network for the cluster to live in:

```shell
# Create network
gcloud compute networks create micro-test --subnet-mode custom --bgp-routing-mode regional

# Create usable subnet
gcloud compute networks subnets create --network=micro-test micro-test --range 192.168.88.0/24 --region us-west2
```

Then create the cluster:

```shell
gcloud container clusters create micro-test --region us-west2 --network micro-test --subnetwork micro-test --num-nodes=1 --machine-type n1-standard-1
```

N.B. In this configuration, num-nodes is number of nodes per availability zone, so you should get 3 nodes in us-west2.


Make sure it's up:

```shell
kubectl config get-contexts
kubectl cluster-info && kubectl get nodes
```

Optionally: Give yourself and members of your team cluster-admin permissions:

```shell
gcloud info | grep -A1 Account
gcloud projects add-iam-policy-binding --member=user:jake@micro.mu --role=roles/container.admin <project id>
```

You need IAM Security Admin permissions over your project to do this.

## Install Micro Prerequisites

Install helm, e.g.

```shell
brew install kubernetes-helm # Or https://github.com/helm/helm#install
kubectl create serviceaccount -n kube-system tiller
kubectl create clusterrolebinding --clusterrole=cluster-admin --serviceaccount=kube-system:tiller tiller-admin
helm init --service-account=tiller --tiller-namespace=kube-system
```

Install etcd

```shell
helm repo update
helm install stable/etcd-operator --version="0.10.0" --set customResources.createEtcdClusterCRD=true --set etcdCluster.version="3.3.15" --set etcdCluster.image.tag="v3.3.15"
```

Install nats

```shell
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml
kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml
kubectl apply -f services/infra/nats.yaml
```

## Install Micro

```shell
grep -rl do-block-storage services/micro | xargs sed -i '' '/do-block/d' # remove first '' for GNU sed
# check services/micro/secrets.yaml
kubectl apply -f services/micro
```

# TODO:

Work out what the firewall rules are to get the micro network running:

```bash
# Turn off the firewall
gcloud compute firewall-rules create inter-instance-communication --network micro-test --allow tcp,udp,icmp --source-ranges 0.0.0.0/0
```
