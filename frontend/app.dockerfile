FROM node:fermium-alpine3.14

WORKDIR /app

EXPOSE 8080
CMD ["npm", "run", "start:dev"]