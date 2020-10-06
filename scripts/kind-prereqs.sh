# this script installs the prerequisites to get kind install working
mkdir /tmp/cfssl
cd /tmp/cfssl
git clone https://github.com/cloudflare/cfssl.git
pushd cfssl 
make
popd
echo "::add-path::$(pwd)/cfssl/bin"
# yq is used to manipulate yaml
GO111MODULE=on go get github.com/mikefarah/yq/v3
