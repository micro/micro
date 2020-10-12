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

MOD=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/micro/v3/proto/api

# modify the proto paths
function modify_paths() {
	while read line; do
		while read path; do 
			m=`echo $path | sed -e 's@go-micro/@go-micro/v2/@g' -e 's@/[a-z]*\.proto$@@g'`
			MOD="${MOD},M${path}=$m"
		done < <(grep "github.com/micro/go-micro/.*\.proto" $line | sed -e 's/^import //g' -e 's/;$//g' -e 's/"//g')
	done < <(find . -name "*.proto")
}

modify_paths
echo Modifiers used: $MOD
##Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/micro/v3/proto/api

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
