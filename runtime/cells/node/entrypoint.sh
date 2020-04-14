#!/bin/dumb-init /bin/sh

set -x  
set -e

REPO=$1

mkdir /app
cd /app

# clone the repo
echo "Cloning $REPO"
git clone $REPO .

# If parameter 2nd parameter is supplied, it's the path
if [ $# -eq 2 ]
  then
    cd $2
fi

npm install
node index.js