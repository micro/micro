FROM alpine:latest AS builder
RUN apk --no-cache add make git go gcc libtool musl-dev

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /
RUN make ; rm -rf $GOPATH/pkg/mod

RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

ENTRYPOINT ["/micro"]