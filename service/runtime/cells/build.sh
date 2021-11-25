#!/bin/bash

set -e

if [ ! $IMAGE ]; then
  IMAGE=micro/cells
fi

echo ${PASSWORD} | docker login $DOCKER_DOMAIN -u ${USERNAME} --password-stdin


ls | while read dir; do
  if [ ! -d ${dir} ]; then
    continue
  fi
  if [ $DOCKER_DOMAIN ]; then
    TAGPREFIX=$DOCKER_DOMAIN/
  fi
  TAG=$TAGPREFIX$IMAGE:${dir}

  pushd ${dir} &>/dev/null
  echo Building $TAG

  if [ ! -s Dockerfile ]; then
    echo Skipping $TAG
    popd &>/dev/null
    continue
  fi

  docker buildx build --platform linux/amd64 --platform linux/arm64 --tag $TAG --push .

  popd &>/dev/null
done
