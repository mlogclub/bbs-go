FROM node:20 AS builder

ENV APP_HOME /code/bbs-go/admin
WORKDIR "$APP_HOME"

COPY . ./
# RUN npm install -g pnpm --registry=https://registry.npmmirror.com
# RUN pnpm install --registry=https://registry.npmmirror.com
RUN npm install -g pnpm
RUN pnpm install
RUN pnpm build:docker

FROM nginx:latest

COPY --from=builder /code/bbs-go/admin/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
