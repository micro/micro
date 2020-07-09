#!/bin/dumb-init /bin/sh

set -x  
set -e

mkdir /app
cd /app

REPO=$(echo $1 | cut -d/ -f -3)
PATH=$(echo $1 | cut -d/ -f 4-)

# clone the repo
echo "Cloning $REPO"
git clone $REPO .

cd $PATH

# run the source
echo "Running service"
go run .