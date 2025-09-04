FROM nginx:latest

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY www /usr/share/nginx/html

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
