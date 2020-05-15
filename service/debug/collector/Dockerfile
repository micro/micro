FROM golang:1.13-alpine AS build-collector

RUN mkdir -p /tmp/build
WORKDIR /tmp/build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN go get github.com/micro/micro/debug/collector@master

FROM netdata/netdata:latest-amd64
COPY --from=build-collector --chown=root:netdata /go/bin/collector /usr/libexec/netdata/plugins.d/micro.d.plugin
