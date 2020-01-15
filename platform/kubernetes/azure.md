# Setting up Micro on Azure

## Prerequisites 

Install the Azure cli using your favourite method, e.g.
```shell
brew install azure-cli
```

Log in to Azure
```shell
az configure
az login # --use-device-code if in a remote shell
```

Install aks-engine
```shell
brew install Azure/aks-engine/aks-engine
```

Install `kubectl`, e.g.
```shell
brew install kubernetes-cli
```

## Create a AKS Engine cluster

Everything in Azure is contained in resource groups, we need one now:

```shell
# Find your preferred location
az account list-locations

# Create resource group
az group create --name micro-cloud-k8s --location eastasia
```

This will output an ID like
`/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/micro-cloud-k8s`
that you need for the next step


Create a Service Principal for the cluster, so it can manage resources in the resource group
```shell
az ad sp create-for-rbac --role="Contributor" --scopes="<id from previous step>"
```

Again, make a note of the appID and Password that you get back from the azure cli.

Then create the cluster

```shell
export SUBSCRIPTION_ID=<subscription uuid>
export CLIENT_ID=<appID>
export CLIENT_SECRET=<password>
aks-engine deploy --subscription-id $SUBSCRIPTION_ID \
    --resource-group micro-cloud-k8s \
    --dns-prefix micro-ap \
    --location eastasia \
    --api-model azure-aks.json \
    --client-id "$CLIENT_ID" \
    --client-secret "$CLIENT_SECRET" \
    --set servicePrincipalProfile.clientID "$CLIENT_ID" \
    --set servicePrincipalProfile.secret "$CLIENT_SECRET" \
    --force-overwrite
```

The kubeconfig and ssh keys will be available in `_output/kubeconfig`

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
kubectl apply -f services/micro
```
