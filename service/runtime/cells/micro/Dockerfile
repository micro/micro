FROM micro/micro
FROM golang:1.15
 
WORKDIR /
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git && \
    git version

COPY entrypoint.sh /
RUN chmod 755 entrypoint.sh
ENTRYPOINT ["bash", "/entrypoint.sh"]
