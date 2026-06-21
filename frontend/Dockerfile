FROM node:22-alpine AS build
WORKDIR /app
ARG PUBLIC_API_URL=http://localhost:16125
ENV PUBLIC_API_URL=$PUBLIC_API_URL
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM docker.io/library/nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
