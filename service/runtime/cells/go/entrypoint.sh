#!/bin/dumb-init /bin/sh

set -x  
set -e

git version

mkdir /app
cd app

REPO=$(echo $1 | cut -d/ -f -3)
P=$(echo $1 | cut -d/ -f 4-)

# clone the repo
echo "Cloning $REPO"
git clone https://$REPO .

cd $P

# run the source
echo "Running service"
GOPROXY=direct go run .