#!/bin/bash
MONITORING_NAMESPACE="monitoring"

# Make sure we have a "monitoring" namespace:
kubectl create namespace ${MONITORING_NAMESPACE}

# Make sure we have the stable repo:
helm repo add stable https://kubernetes-charts.storage.googleapis.com

# Install Prometheus using Helm:
helm install prometheus stable/prometheus \
    --namespace ${MONITORING_NAMESPACE} \
    --set alertmanager.enabled=false \
    --set pushgateway.enabled=false
