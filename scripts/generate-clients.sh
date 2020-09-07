#!/bin/sh -e

### UBUNTU BIONIC ###
echo "deb https://packages.le-vert.net/tensorflow/ubuntu bionic main" | sudo tee -a /etc/apt/sources.list
wget -O - https://packages.le-vert.net/packages.le-vert.net.gpg.key | sudo apt-key add -

sudo apt update
sudo apt install -y protobuf-compiler
sudo apt install -y --no-install-recommends python3 python3-pip python3-setuptools python3-dev python3-grpcio python3-protobuf
sudo apt install -y --no-install-recommends nodejs npm
sudo apt install -y --no-install-recommends ruby
sudo gem install grpc grpc-tools
pip3 install --no-cache-dir grpcio-tools
npm i grpc-tools
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
cargo install protobuf-codegen grpc-compiler
go mod init clients
GOBIN=$(go env GOPATH)/bin go get github.com/micro/micro/v3/cmd/protoc-gen-client@master github.com/golang/protobuf/protoc-gen-go@v1.4.1 github.com/micro/micro/v3/cmd/protoc-gen-micro github.com/micro/go-micro/v3@latest
