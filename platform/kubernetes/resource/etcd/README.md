# Etcd

Etcd is managed through helm

## Usage

To start etcd

```
./install.sh
```

To delete etcd

```
# list existing deployments
helm list

# remove the deployment
helm delete etcd-cluster
```

