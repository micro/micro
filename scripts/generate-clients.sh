#!/bin/bash

set -x
set -e

GO_PATH=$(go env GOPATH)
GO_BIN=$GO_PATH/bin
GO=$GO_PATH/bin/go
PATH=$PATH:$GO_BIN:$(npm bin):/usr/local/bin/:$HOME/.cargo/bin

### UBUNTU BIONIC ###
echo "deb https://packages.le-vert.net/tensorflow/ubuntu bionic main" | sudo tee -a /etc/apt/sources.list
wget -O - https://packages.le-vert.net/packages.le-vert.net.gpg.key | sudo apt-key add -

# install all the deps
sudo apt update
sudo apt install -y protobuf-compiler
sudo apt install -y --no-install-recommends python3 python3-pip python3-setuptools python3-dev python3-grpcio python3-protobuf
sudo apt install -y --no-install-recommends nodejs npm
sudo apt install -y --no-install-recommends ruby ruby-dev
sudo apt install -y --no-install-recommends git-all
sudo gem update --system
sudo gem install grpc grpc-tools
pip3 install --no-cache-dir grpcio-tools
npm i grpc-tools
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
cargo install protobuf-codegen grpc-compiler

# build proto related code
pushd cmd/protoc-gen-client && go get ./... && popd
pushd cmd/protoc-gen-micro && go get ./... && popd
go get github.com/golang/protobuf/protoc-gen-go@v1.4.2
# delete the existing sdk directory
rm -rf client/sdk
# generate the clients
PATH=$PATH:$GO_BIN:$(npm bin):/usr/local/bin/:$HOME/.cargo/bin protoc-gen-client -srcdir proto/ -dstdir client/sdk/ -langs go,python,java,ruby,node,rust
# remove node garbage
rm -rf node_modules/ package-lock.json
