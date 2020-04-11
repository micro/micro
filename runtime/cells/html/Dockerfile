FROM nginx:stable-alpine

RUN apk add --no-cache git
RUN sed -i '1idaemon off;' /etc/nginx/nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY entrypoint.sh /
RUN chmod 755 /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["sh", "/entrypoint.sh"]