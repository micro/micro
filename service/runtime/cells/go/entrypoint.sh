#!/bin/dumb-init /bin/sh

set -x
set -e

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
git clone $CLONE_URL --branch $REF --single-branch .
if [ $? -eq 0 ]; then
    echo "Successfully cloned branch"
else
    # Clone the full repo if the REF was not a branch.
    # In case of a commit REF we will git reset later.
    git clone https://$REPO .
fi

# Try to check out commit and do not care if it fails
git reset --hard $REF

cd $P

# find the entrypoint using the util
ENTRYPOINT=$(entrypoint)

# run the source
echo "Running service"
go run $ENTRYPOINT