#!/bin/bash

set -e

function trap_handler {
  MYSELF="$0"   # equals to my script name
  LASTLINE="$1" # argument 1: last line of error occurence
  LASTERR="$2"  # argument 2: error code of last command
  echo "Error: line ${LASTLINE} - exit status of last command: ${LASTERR}"
  exit $2
}
trap 'trap_handler ${LINENO} ${$?}' ERR

echo "Checking dependencies..."
which protoc
which protoc-gen-go

echo "Building protobuf code..."
DIR=`pwd`
SRCDIR=`cd $DIR && cd ../../.. && pwd`
find $DIR/proto -name '*.pb.go' -exec rm {} \;
find $DIR/proto -name '*.proto' -exec echo {} \;
find $DIR/proto -name '*.proto' -exec protoc -I$SRCDIR --go_out=plugins=micro:${SRCDIR} {} \;

echo "Complete"
