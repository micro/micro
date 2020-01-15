# Setting up Micro on Amazon Web Services

## Prerequisites 

Install the AWS CLI tools using your preferred method, e.g.:

```shell
brew install awscli
# pip install awscli --upgrade --user
```

Grab your access key ID and Secret Access Key, then configure. 

```shell
aws configure
AWS Access Key ID [********************]:
AWS Secret Access Key [********************]:
Default region name [ap-east-1]:
Default output format [json]:
```

Make sure you have kubectl installed, e.g.

```shell
brew install kubernetes-cli
```

The recommended way to use EKS is to use [eksctl](https://eksctl.io/), so install that:

```shell
# See https://github.com/weaveworks/eksctl#Installation
curl | sudo bash
eksctl version
[ℹ]  version.Info{BuiltAt:"", GitCommit:"", GitTag:"0.6.0"}
```

In order to log to any cluster you create you will also need aws-iam-authenticator:

```shell
# https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html
brew install aws-iam-authenticator
```

## Create an EKS Cluster

Create a config file that represents the cluster you want:

```shell
cat aws-eks.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: micro
  region: ap-east-1
  version: "1.14"

nodeGroups:
  - name: micro-workers-1
    instanceType: m5.large
    desiredCapacity: 3
    iam:
      withAddonPolicies:
        ebs: true
```

```shell
eksctl create cluster -f aws-eks.yaml
```

This takes a while, go and get coffee. If it failed for some reason delete the stack with `eksctl delete cluster --name=micro --region=ap-east-1` (Or whatever you tried to create). Then check the AWS console for cloudformation stacks that might have been left behind and delete those too.

Once you're successfully up, go to the ClusterSharedNodeSecurityGroup in EC2 -> Security groups and allow all UDP traffic in (TODO: lock down to only the ports we require.)

You can also grant other people access to your cluster.

Edit the configmap `aws-auth` in the `kube-system` namespace:

```shell
kubectl edit configmap -n kube-system aws-auth
```

add people to the `mapUsers:` block:

```yaml
  mapUsers: |
    - userarn: arn:aws:iam::555555555555:user/someone
      username: someone
      groups:
      - system:masters
```

They can then log in with the aws-cli:

```shell
aws eks get-token --cluster-name mycluster
```

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

## Deploy micro

```shell
grep -rl do-block-storage services/micro | xargs sed -i '' '/do-block/d' # remove first '' for GNU sed
# check services/micro/secrets.yaml
kubectl apply -f services/micro
```

