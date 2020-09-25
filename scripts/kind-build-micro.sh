#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

if [[ ! -d $TMPDIR/micro-kind ]]; then
  mkdir $TMPDIR/micro-kind
fi
rsync -av --exclude=$DIR/platform/kubernetes $DIR/../* $TMPDIR/micro-kind/

pushd $TMPDIR/micro-kind
micro init --profile=platform --output=profile.go
go mod edit -replace github.com/micro/micro/profile/platform/v3=./profile/platform
go mod edit -replace google.golang.org/grpc=google.golang.org/grpc@v1.26.0
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
docker build -t micro -f test/Dockerfile-kind .
docker tag micro localhost:5000/micro
docker push localhost:5000/micro
popd