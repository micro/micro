FROM alpine
RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*
ADD micro /micro
ADD web/webapp/dist /web/webapp/dist
WORKDIR /
ENTRYPOINT [ "/micro" ]
