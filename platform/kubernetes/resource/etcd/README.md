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

Note: When connecting to etcd, the ca and client certs must be used, e.g:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: micro-registry
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: registry
        env:
        - name: MICRO_REGISTRY
          value: "etcd"
        - name: MICRO_REGISTRY_ADDRESS
          value: "etcd-cluster"
        - name: MICRO_REGISTRY_SECURE
          value: "true"
        - name: MICRO_CERTIFICATE_AUTHORITIES
          value: "/certs/registry/ca.crt"
        args:
        - registry
        image: bentoogood/micro:mtls
        imagePullPolicy: Always
        ports:
        - containerPort: 8000
          name: registry-port
        volumeMounts:
        - name: etcd-client-certs
          mountPath: "/certs/registry"
          readOnly: true
      volumes:
      - name: etcd-client-certs
        secret:
          secretName: etcd-client-certs
```