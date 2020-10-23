FROM golang:1.13.15-alpine3.12 AS builder
WORKDIR /micro
RUN apk --no-cache add make git gcc libtool musl-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make ; rm -rf $GOPATH/pkg/mod

FROM golang:1.13.15-alpine3.12
RUN apk --no-cache add make git gcc libtool musl-dev
RUN apk --no-cache add ca-certificates && rm -rf /var/cache/apk/* /tmp/* 
COPY --from=builder /micro/micro /micro
ENTRYPOINT ["/micro"]