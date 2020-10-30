#!/bin/bash

helm uninstall nginx
kubectl delete namespace cert-manager

kubectl delete -f letsencrypt.yaml
kubectl delete -f ingress.yaml
kubectl delete secret nginx-tls