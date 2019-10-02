# Kubernetes Deployment

## Dependencies

- Kubectl
- Kubectx
- Helm
- Etcd
- NATS

## Steps

1. Spin up managed k8s on DO/GCP/AWS

2. Setup local env
  - Install `kubectl` https://kubernetes.io/docs/tasks/tools/install-kubectl/
  - Install `kubectx` https://github.com/ahmetb/kubectx
  - Install `helm` https://github.com/helm/helm
    * See [this](https://github.com/helm/helm/blob/master/docs/rbac.md)

3. Install `etcd` on DO/GCP/AWS
  - `helm repo update`
  - `helm install stable/etcd-operator --version="0.10.0" --set customResources.createEtcdClusterCRD=true`
  - read the docs [here](https://etcd.io/docs/v3.2.17/)

4. Install Micro core on DO/GCP/AWS
  - kubectl apply -f ../kubernetes
  - Create external load balancers https://www.digitalocean.com/docs/kubernetes/how-to/add-load-balancers/

5. Install `etcd` in your local environment and query remote `etcd` cluster
  - `kubectl port-forward service/etcd-cluster-client -n default 2379`
  - `ETCDCTL_API=3 etcdctl version`
  - `ETCDCTL_API=3 etcdctl -w table member list`

6. Accessing particular `etcd` node directly
  - `kubectl exec -it ETCD_NODE_POD_NAME -- /bin/sh`
