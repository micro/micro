#!/bin/bash

set -x

## Submit an update for go mod

sed -i 's@github.com/micro/go-micro/v2 .*@github.com/micro/go-micro/v2 master@g' go.mod
go fmt
git add go.mod
git branch -D go-mod
git branch go-mod
git checkout go-mod
git add go.mod go.sum
git commit -m "Update go.mod"
git push origin go-mod
git checkout master
git branch -D go-mod
