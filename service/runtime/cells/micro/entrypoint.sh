#!/bin/dumb-init /bin/sh

set -x  

git version

mkdir /app
cd app

URL=$1

REF=""
if [[ $URL == *"@"* ]]; then
   # Save the ref
  REF=$(echo $URL | cut -d@ -f 2-)
  # URL should not contain the ref
  URL=$(echo $URL | cut -d@ -f -1)
fi

if [[ $REF == "latest" ]]; then
  REF="master"
fi

REPO=$(echo $URL | cut -d/ -f -3)
P=$(echo $URL | cut -d/ -f 4-)

echo "Repo is $REPO"
echo "Path is $P"
echo "Ref is $REF"


if [[ -z "$GIT_CREDENTIALS" ]]; then 
  echo "Cloning $REPO"
  CLONE_URL=https://$REPO
else
  echo "Cloning $REPO with credentials"
  CLONE_URL=https://$GIT_CREDENTIALS@$REPO
fi

# clone the repo
git clone $CLONE_URL --depth=1
cd $P

git fetch origin $REF --depth 1

git checkout $REF

# run the source
echo "Running service"
go run .