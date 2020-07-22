FROM alpine:latest
RUN apk --no-cache add make git go gcc libtool musl-dev

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

COPY . /
RUN make

FROM alpine:latest
# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin # not sure if we need this
RUN apk --no-cache add ca-certificates \
    gcc \
    git \
    go \
    libtool \
    musl-dev \
    && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=0 /micro /micro
ENTRYPOINT ["/micro"]
