# This is mostly intended to be triggered by CI
# as it modifies the source code.

git clone https://github.com/cloudflare/cfssl.git
pushd cfssl 
make
popd

GO111MODULE=on go get github.com/mikefarah/yq/v3

yq write -i platform/kubernetes/network/proxy.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" false
yq write -i platform/kubernetes/network/router.yaml "spec.template.spec.containers[0].env.(name==MICRO_ENABLE_ACME).value" false
yq delete -i platform/kubernetes/network/proxy.yaml "spec.template.spec.containers[0].env.(name==CF_API_TOKEN)"

pushd platform/kubernetes
./install.sh
popd

# TODO 
# how do we make it pull down this version of micro ?
# sed -i "/latest/ s/latest/$GITHUB_BRANCH/g" micro.tf
