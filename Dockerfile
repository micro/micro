FROM alpine:3.12.1 as builder

COPY --from=golang:1.15-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
RUN apk --no-cache add make git gcc libtool musl-dev

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /
RUN make ; rm -rf $GOPATH/pkg/mod

FROM alpine:3.12.1
COPY --from=golang:1.15-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

RUN apk --no-cache add make git gcc libtool musl-dev
RUN apk --no-cache add ca-certificates && rm -rf /var/cache/apk/* /tmp/* 

COPY --from=builder /micro /micro
ENTRYPOINT ["/micro"]