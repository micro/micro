FROM golang:1.13-alpine as builder
RUN apk --no-cache add make git gcc libtool musl-dev
WORKDIR /
COPY . /
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*
COPY --from=builder /micro .
ENTRYPOINT ["/micro"]
