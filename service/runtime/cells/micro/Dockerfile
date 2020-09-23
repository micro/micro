FROM alpine:latest AS builder
RUN apk --no-cache add make git go gcc libtool musl-dev
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR loader
COPY loader .
RUN go build  .

FROM alpine:latest
COPY --from=builder loader/loader loader
CMD ./loader