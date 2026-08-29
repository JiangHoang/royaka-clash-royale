# Stage 1: Build
FROM node:18 AS build
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
ARG VITE_WS_URL=ws://localhost:8080/ws
ENV VITE_WS_URL=${VITE_WS_URL}
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html/royaka-2025-fe
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
