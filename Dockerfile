FROM golang:1.15-alpine3.12 AS builder
RUN apk --no-cache add make git gcc libtool musl-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /
RUN make ; rm -rf $GOPATH/pkg/mod

FROM golang:1.15-alpine3.12
RUN apk --no-cache add make git gcc libtool musl-dev
RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=builder /micro /micro
ENTRYPOINT ["/micro"]