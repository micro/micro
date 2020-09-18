#!/bin/bash
MONITORING_NAMESPACE="monitoring"

# Make sure we have a "monitoring" namespace:
kubectl create namespace ${MONITORING_NAMESPACE}

# Make sure we have the stable repo:
helm repo add stable https://kubernetes-charts.storage.googleapis.com

# Install Prometheus using Helm:
helm upgrade prometheus stable/prometheus \
    --install \
    --namespace ${MONITORING_NAMESPACE} \
    --set alertmanager.enabled=false \
    --set alertmanager.persistentVolume.enabled=false \
    --set pushgateway.enabled=false \
    --set server.persistentVolume.enabled=false \
    --set-file extraScrapeConfigs=extraScrapeConfigs.yaml
