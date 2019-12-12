FROM golang:1.13 AS buildplugin
ENV GO111MODULE on
RUN mkdir /tmp/build
WORKDIR /tmp/build
RUN git clone https://github.com/jakexks/netdata-collector
WORKDIR /tmp/build/netdata-collector
ENV CGO_ENABLED 0
RUN go build .

FROM netdata/netdata:latest
COPY --from=buildplugin /tmp/build/netdata-collector/netdata-collector /usr/libexec/netdata/plugins.d/micro.d.plugin
RUN chown root:netdata /usr/libexec/netdata/plugins.d/micro.d.plugin
RUN chmod 0750 /usr/libexec/netdata/plugins.d/micro.d.plugin
