# Etcd

Etcd is managed through helm

## Usage

To start etcd

```
./install.sh
```

To delete etcd

```
./uninstall.sh
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
        - name: MICRO_REGISTRY_TLS_CA
          value: "/certs/registry/ca.crt"
        - name: MICRO_REGISTRY_TLS_CERT
          value: "/certs/registry/cert.pem"
        - name: MICRO_REGISTRY_TLS_KEY
          value: "/certs/registry/key.pem"
        args:
        - registry
        image: micro/micro:mtls
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