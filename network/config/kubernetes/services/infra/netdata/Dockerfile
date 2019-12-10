FROM golang:1.13 AS buildplugin
ENV GO111MODULE on
RUN mkdir /tmp/build
WORKDIR /tmp/build
RUN git clone https://github.com/jakexks/go.d.plugin
WORKDIR /tmp/build/go.d.plugin
RUN git checkout gomicro
RUN go build ./cmd/godplugin

FROM netdata/netdata:latest
COPY --from=buildplugin /tmp/build/go.d.plugin/godplugin /usr/libexec/netdata/plugins.d/go.d.plugin
RUN chown root:netdata /usr/libexec/netdata/plugins.d/go.d.plugin
RUN chmod 0750 /usr/libexec/netdata/plugins.d/go.d.plugin
