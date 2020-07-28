FROM nginx:1.17.10-alpine AS fileserver

COPY nginx.conf /etc/nginx/conf.d/default.conf
ADD images /web/www/static/

EXPOSE 80