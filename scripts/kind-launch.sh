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
# safety first
kubectl config use-context kind-kind

# yq does not quote the strings in it's output so need to hack this and fix with sed so it doesn't get parsed as a bool later in k8s
yq write -i platform/kubernetes/network/router.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" '_false_'
sed -e 's/_false_/"false"/g' -i.bak platform/kubernetes/network/router.yaml
yq write -i platform/kubernetes/network/proxy.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" '_false_'
sed -e 's/_false_/"false"/g' -i.bak platform/kubernetes/network/proxy.yaml
yq delete -i platform/kubernetes/network/proxy.yaml "spec.template.spec.containers[0].env.(name==CF_API_TOKEN)"

pushd platform/kubernetes
./install.sh
popd

# TODO 
# how do we make it pull down this version of micro ?
# Build and push to docker a named version
# update the yamls