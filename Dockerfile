FROM alpine:3.2
ADD micro /micro
WORKDIR /
ENTRYPOINT [ "/micro" ]
