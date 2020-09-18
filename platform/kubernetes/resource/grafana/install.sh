#!/bin/bash
MONITORING_NAMESPACE="monitoring"

# Make sure we have a "monitoring" namespace:
kubectl create namespace ${MONITORING_NAMESPACE}

# Make sure we have the stable repo:
helm repo add stable https://kubernetes-charts.storage.googleapis.com

# Install Grafana using Helm:
helm upgrade grafana stable/grafana \
    --install \
    --namespace ${MONITORING_NAMESPACE} \
    --set persistence.enabled=true \
    --set persistence.size=1Gi
