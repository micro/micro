# Kubernetes Deployment

## Dependencies

- Kubectl
- Kubectx
- Helm
- Consul
- NATS

## Steps

1. Spin up managed k8s on DO, GCP, AWS

2. Setup local env
  - Install Kubectl https://kubernetes.io/docs/tasks/tools/install-kubectl/
  - Install Kubectx https://github.com/ahmetb/kubectx
  - Install Helm https://github.com/helm/helm
    * See https://github.com/helm/helm/blob/master/docs/rbac.md
    
3. Install Consul
  - https://www.consul.io/docs/platform/k8s/run.html
  - kubectl port-forward consul-consul-server-0 8500:8500
  
4. Install Micro
  - kubectl apply -f ../kubernetes  
  - Create external load balancers https://www.digitalocean.com/docs/kubernetes/how-to/add-load-balancers/
                         
