#!/bin/bash

IMAGE=micro/cells

echo ${PASSWORD} | docker login ghcr.io -u ${USERNAME} --password-stdin


ls | while read dir; do
  if [ ! -d ${dir} ]; then
    continue
  fi

  TAG=ghcr.io/micro/$IMAGE:${dir}

  pushd ${dir} &>/dev/null
  echo Building $TAG

  if [ ! -s Dockerfile ]; then
    echo Skipping $TAG
    popd &>/dev/null
    continue
  fi

  docker build -t $TAG .
  docker push $TAG

  popd &>/dev/null
done
