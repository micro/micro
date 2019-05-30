#!/bin/bash

set -e
#set -x

# Used to test all the things

find * -type d -maxdepth 0 -print | while read dir; do
	if [ -f $dir/.skip ]; then
		continue
	fi

	pushd $dir >/dev/null

	echo 
	echo "Testing $dir"
	echo

	# run tests
	go test -v ./...

	popd >/dev/null
done
