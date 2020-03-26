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

MOD=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto

echo "Building protobuf code..."
DIR=`pwd`
SRCDIR=$GOPATH/src

echo "DIR $DIR"
echo "SRCDIR $SRCDIR"
find $DIR/proto -name '*.pb.go' -exec rm {} \;
find $DIR/proto -name '*.micro.go' -exec rm {} \;
find $DIR/proto -name '*.proto' -exec echo {} \;
find $DIR/proto -name '*.proto' -exec protoc --proto_path=$SRCDIR --micro_out=${MOD}:${SRCDIR} --go_out=${MOD}:${SRCDIR} {} \;
#find $DIR/proto -name '*.proto' -exec protoc --proto_path=$SRCDIR --micro_out=${SRCDIR} --go_out=${SRCDIR} {} \;
#find $DIR/proto -name '*.proto' -exec protoc --proto_path=$SRCDIR --micro_out=${SRCDIR}:${MOD} --go_out=plugins=grpc:${SRCDIR} {} \;

echo "Complete"
