# this script installs the prerequisites to get kind install working
mkdir /tmp/cfssl
cd /tmp/cfssl
git clone https://github.com/cloudflare/cfssl.git
pushd cfssl 
make
popd
echo "$(pwd)/cfssl/bin" >> $GITHUB_PATH
# yq is used to manipulate yaml
GO111MODULE=on go get github.com/mikefarah/yq/v3
