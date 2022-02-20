FROM node:fermium-alpine3.14

RUN apk update \
    && apk add wget bash \
    && wget -O /opt/wait-for-it.sh -q https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
    && chmod -R 755 /opt/wait-for-it.sh

WORKDIR /app

EXPOSE 3000
ENTRYPOINT [ "/bin/sh", "entrypoint.sh" ]
CMD ["npm","run","start:dev"]