FROM micro/go-micro
FROM micro/micro

WORKDIR /
COPY entrypoint.sh /
RUN chmod 755 entrypoint.sh
ENTRYPOINT ["bash", "/entrypoint.sh"]
