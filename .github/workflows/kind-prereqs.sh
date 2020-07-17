# this script installs the prerequisites to get kind install working
git clone https://github.com/cloudflare/cfssl.git
pushd cfssl 
make
popd
echo "::addpath::$(pwd)/cfssl/bin"
# yq is used to manipulate yaml
GO111MODULE=on go get github.com/mikefarah/yq/v3
