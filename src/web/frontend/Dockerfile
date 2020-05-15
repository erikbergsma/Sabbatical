FROM nginx

EXPOSE 80/tcp

COPY ./nginx.conf.template /tmp/nginx.conf.template
COPY ./webroot/* /usr/share/nginx/html/

ENTRYPOINT /bin/bash -c "envsubst < /tmp/nginx.conf.template > /etc/nginx/nginx.conf && nginx -g 'daemon off;'"
