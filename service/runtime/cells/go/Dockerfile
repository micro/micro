FROM golang:1.15.3

WORKDIR /
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git && \
    git version

# install the entrypoint helper (finds the main.go)
COPY util util
WORKDIR util/entrypoint
RUN go install

WORKDIR /
COPY entrypoint.sh /
RUN chmod 755 entrypoint.sh
ENTRYPOINT ["bash", "/entrypoint.sh"]
