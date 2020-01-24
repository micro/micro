FROM golang:1.13-alpine as builder
RUN apk --no-cache add make git gcc libtool musl-dev upx
WORKDIR /
COPY . /
RUN make build && upx micro

FROM alpine:latest

RUN apk add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=builder /micro .
ENTRYPOINT ["/micro"]
