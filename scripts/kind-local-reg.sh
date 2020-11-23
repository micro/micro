#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

if [[ ! -d $TMPDIR/micro-kind ]]; then
  mkdir $TMPDIR/micro-kind
fi

# start with a clean dir
rm -rf $TMPDIR/micro-kind/*
cp -R $DIR/../* $TMPDIR/micro-kind/

pushd $TMPDIR/micro-kind
docker run -d -p 5000:5000 --restart=always --name kind-registry -v /tmp/docker-registry:/var/lib/registry registry:2
./scripts/kind-build-micro.sh

kind create cluster --config ./scripts/kind/kind-config.yaml
docker network connect "kind" "kind-registry"

for node in $(kind get nodes);
do
  kubectl annotate node "${node}" "kind.x-k8s.io/registry=localhost:5000"
done

sed_expression="s/: ghcr.io\/m3o\/platform/: localhost:5000\/micro/g"
sed -e "$sed_expression" -i.bak ./cmd/platform/kubernetes/service/*.yaml

popd
