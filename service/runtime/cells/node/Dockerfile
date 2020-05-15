FROM node:10

COPY entrypoint.sh /
RUN chmod 755 entrypoint.sh

EXPOSE 8080
ENTRYPOINT [ "sh", "/entrypoint.sh" ]