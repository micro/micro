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
    cd $2
fi

# run the source
echo "Running service"
go run .