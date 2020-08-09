#!/bin/bash

GOOS=linux CGO_ENABLED=0 GO111MODULE=on go build -a -installsuffix cgo -ldflags '-w' -o collector main.go
