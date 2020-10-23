FROM golang:1.15-alpine3.12
RUN apk add make git gcc libtool musl-dev curl bash
RUN apk add ca-certificates && rm -rf /var/cache/apk/* /tmp/* 

RUN         apk add --update ca-certificates openssl tar && \
            wget https://github.com/coreos/etcd/releases/download/v3.4.7/etcd-v3.4.7-linux-amd64.tar.gz && \
            tar xzvf etcd-v3.4.7-linux-amd64.tar.gz && \
            mv etcd-v3.4.7-linux-amd64/etcd* /bin/ && \
            apk del --purge tar openssl && \
            rm -Rf etcd-v3.4.7-linux-amd64* /var/cache/apk/*
VOLUME      /data
EXPOSE      2379 2380 4001 7001
ADD         scripts/run-etcd.sh /bin/run.sh
COPY . .
ENTRYPOINT ["sh", "/bin/run.sh"]