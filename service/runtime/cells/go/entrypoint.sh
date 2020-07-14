#!/bin/dumb-init /bin/sh

set -x  
set -e

git version

mkdir /app
cd app

URL=$1
if [[ $1 != *"github"* ]]; then
  URL="github.com/micro/services/$URL"
fi

REPO=$(echo $URL | cut -d/ -f -3)
P=$(echo $URL | cut -d/ -f 4-)

echo "Repo is $REPO"
echo "Path is $P"

# clone the repo
echo "Cloning $REPO"
git clone https://$REPO .

cd $P

# run the source
echo "Running service"
go run .