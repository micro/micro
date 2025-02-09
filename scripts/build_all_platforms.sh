#!/bin/bash
set -e
builds=("windows,amd64" "windows,386" "linux,arm64" "linux,386" "linux,amd64" "linux,arm" "darwin,amd64")

for build in ${builds[@]}; do
  IFS=',' read -r -a array <<< "$build"
  echo "building ${array[0]} ${array[1]}"
  # TODO: CGO_ENABLED=1 prevent building windows amd64
  GOOS=${array[0]} GOARCH=${array[1]} CGO_ENABLED=0 go build cmd/micro/main.go
done
