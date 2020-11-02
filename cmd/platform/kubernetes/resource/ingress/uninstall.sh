#!/bin/bash

helm uninstall nginx
kubectl delete namespace cert-manager

kubectl delete -f letsencrypt.yaml
kubectl delete -f ingress.yaml
kubectl delete secret api-tls
kubectl delete secret proxy-tls