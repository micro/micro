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

yq write -i platform/kubernetes/network/proxy.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" '"false"'
yq write -i platform/kubernetes/network/router.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" '"false"'
yq delete -i platform/kubernetes/network/proxy.yaml "spec.template.spec.containers[0].env.(name==CF_API_TOKEN)"

pushd platform/kubernetes
./install.sh
popd

# TODO 
# how do we make it pull down this version of micro ?
