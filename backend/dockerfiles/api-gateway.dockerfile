FROM node:fermium-alpine3.14

WORKDIR /app

EXPOSE 3000
CMD ["npm","run","start:dev"]