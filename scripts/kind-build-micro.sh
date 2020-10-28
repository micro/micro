#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

tmp=$TMPDIR
if [[ ! $tmp ]]; then
  tmp=/tmp
fi

if [[ ! -d $tmp/micro-kind ]]; then
  mkdir $tmp/micro-kind
fi
rsync -av --exclude=$DIR/../cmd/platform/kubernetes $DIR/../* $tmp/micro-kind/

pushd $tmp/micro-kind
micro init --profile=platform --output=profile.go
go mod edit -replace github.com/micro/micro/plugin/etcd/v3=./plugin/etcd
go mod edit -replace github.com/micro/micro/plugin/cockroach/v3=./plugin/cockroach
go mod edit -replace github.com/micro/micro/plugin/prometheus/v3=./plugin/prometheus
go mod edit -replace github.com/micro/micro/plugin/nats/broker/v3=./plugin/nats/broker
go mod edit -replace github.com/micro/micro/plugin/nats/stream/v3=./plugin/nats/stream
go mod edit -replace github.com/micro/micro/profile/platform/v3=./profile/platform
go mod edit -replace google.golang.org/grpc=google.golang.org/grpc@v1.26.0
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
docker build -t micro -f test/Dockerfile-kind .
docker tag micro localhost:5000/micro
docker push localhost:5000/micro
popd
