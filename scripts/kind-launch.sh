#!/bin/bash

# This script installs the platform on a local k8s (kind) cluster
# It assumes:
# - kind is already running and kubectl context is pointed to it
# - the following tools are installed
#   - helm
#   - cfssl - https://github.com/cloudflare/cfssl
#   - yq - https://github.com/mikefarah/yq
# 
# Warning: This script will modify some yaml files so please don't commit the modifications

set -e
set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# safety first
kubectl config use-context kind-kind

tmp=$TMPDIR
if [[ ! $tmp ]]; then
  tmp=/tmp
fi

KUBE_DIR=$tmp/micro-kind/cmd/platform/kubernetes

if [[ ! -d $tmp/micro-kind ]]; then
  mkdir $tmp/micro-kind
  cp -R $DIR/../* $tmp/micro-kind/
fi

pushd $tmp/micro-kind

yq write -i $KUBE_DIR/service/router.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" --tag '!!str' 'false'
yq write -i $KUBE_DIR/service/proxy.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" --tag '!!str' 'false'
yq delete -i $KUBE_DIR/service/proxy.yaml "spec.template.spec.containers[0].env.(name==CF_API_TOKEN)"
yq write -i $KUBE_DIR/service/api.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" --tag '!!str' 'false'
yq delete -i $KUBE_DIR/service/api.yaml "spec.template.spec.containers[0].env.(name==CF_API_TOKEN)"
yq write -i $KUBE_DIR/service/api.yaml "spec.template.spec.containers[0].ports.(name==api-port).containerPort" 8080

# install metrics server
kubectl apply -f scripts/kind/metrics/components.yaml

pushd $KUBE_DIR
./install.sh dev
kubectl wait deployment --all --timeout=180s -n default --for=condition=available
popd

popd
