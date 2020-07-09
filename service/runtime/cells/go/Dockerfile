FROM golang:1.13.8

WORKDIR /
COPY entrypoint.sh /
RUN chmod 755 entrypoint.sh
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git
ENTRYPOINT ["bash", "/entrypoint.sh"]
