#!/bin/bash

# Description: update.sh simply updates a local micro deployment

export GOPATH=~/c
export GO111MODULE=on

# get go-micro
pushd ~/c/src/github.com/micro/
pushd go-micro
git pull
popd

# get micro
pushd micro
git checkout go.sum
git pull
rm go.sum
go get ./...
popd

## pop pop
popd

/etc/init.d/micro restart
