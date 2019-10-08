# Kubernetes Deployment

## Dependencies

- Kubectl
- Kubectx
- Helm
- Etcd
- NATS

## Steps

1. Spin up managed k8s on DO/[GCP](gcloud.md)/AWS

2. Setup local env
  - Install `kubectl` https://kubernetes.io/docs/tasks/tools/install-kubectl/
  - Install `kubectx` https://github.com/ahmetb/kubectx
  - Install `helm` https://github.com/helm/helm
    * See [this](https://github.com/helm/helm/blob/master/docs/rbac.md)

3. Install `etcd` on DO/GCP/AWS
  - `helm repo update`
  - `helm install stable/etcd-operator --version="0.10.0" --set customResources.createEtcdClusterCRD=true --set etcdCluster.version="3.3.15" --set etcdCluster.image.tag="v3.3.15"`
  - read the docs [here](https://etcd.io/docs/v3.3.12/)

4. Install `nats` on DO/GCP/AWS
 - `kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml`
 - `kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml`
 - `kubectl apply -f services/infra/nats.yaml`

4. Install Micro core on DO/GCP/AWS
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
