FROM node:20-alpine AS builder

WORKDIR /ui

COPY package*.json ./

RUN npm install

COPY ./src ./src
COPY ./public ./public
COPY ./index.html ./
COPY ./svelte.config.js ./
COPY ./vite.config.ts ./
COPY ./tsconfig*.json ./
COPY ./components.json ./
COPY ./postcss.config.js ./
COPY ./tailwind.config.ts ./
COPY ./package-lock.json ./
COPY ./package.json ./

RUN npm run build

FROM nginx:alpine

COPY --from=builder /ui/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]