#!/bin/bash

set -x

## Submit an update for go mod
git branch -D go-mod
git branch go-mod
git checkout go-mod
git pull origin master
GOPROXY=direct go get github.com/micro/go-micro/v3@master
go fmt
go mod tidy
git add go.mod go.sum
git commit -m "Update go.mod"
git push origin go-mod
git checkout master
git branch -D go-mod
