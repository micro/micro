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

DIR=cmd/platform/kubernetes

# safety first
kubectl config use-context kind-kind

yq write -i $DIR/service/router.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" --tag '!!str' 'false'
yq write -i $DIR/service/proxy.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" --tag '!!str' 'false'
yq delete -i $DIR/service/proxy.yaml "spec.template.spec.containers[0].env.(name==CF_API_TOKEN)"
yq write -i $DIR/service/api.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" --tag '!!str' 'false'
yq delete -i $DIR/service/api.yaml "spec.template.spec.containers[0].env.(name==CF_API_TOKEN)"
yq write -i $DIR/service/api.yaml "spec.template.spec.containers[0].ports.(name==api-port).containerPort" 8080

# install metrics server
kubectl apply -f scripts/kind/metrics/components.yaml

pushd $DIR
./install.sh dev
kubectl wait deployment --all --timeout=180s -n default --for=condition=available 
popd
