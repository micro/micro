FROM alpine:latest
RUN apk add make git go gcc libtool musl-dev curl bash

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN apk add avahi avahi-tools

RUN echo -e '#!/bin/bash\navahi-daemon &' > run.sh && chmod +x run.sh && ./run.sh

RUN apk add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
ADD         scripts/run-etcd.sh /bin/run.sh

# Speed up tests by predownloading and building dependencies for services used.
COPY . .
RUN go get github.com/micro/services
RUN go get github.com/micro/services/helloworld
RUN bash -c 'for d in $(find test -name "main.go" | xargs -n 1 dirname); do pushd $d && go get && popd; done'
COPY ./micro /micro
ENTRYPOINT ["/micro"]
