#!/bin/dumb-init /bin/sh

set -x  
set -e

mkdir /app
cd /app

REPO=$1

# clone the repo
echo "Cloning $REPO"
git clone $REPO .

# If parameter 2nd parameter is supplied, it's the path
if [ $# -eq 2 ]
  then
    echo "cding into $2"
    mv $2/* /app/
fi

echo "Serving single page application"
nginx