# Etcd

Etcd is managed through helm

## Usage

To start etcd

```
helm install stable/etcd-operator --version="0.10.0" --set customResources.createEtcdClusterCRD=true --set etcdCluster.version="3.4.3" --set etcdCluster.image.tag="v3.4.3"
```

To delete etcd

```
# list existing deployments
helm list

# remove the deployment
helm delete [name]
```
