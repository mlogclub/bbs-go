FROM node:20-alpine AS builder

ENV APP_HOME /code/bbs-go/site
WORKDIR "$APP_HOME"

COPY . ./
# RUN npm install -g pnpm --registry=https://registry.npmmirror.com
# RUN pnpm install --registry=https://registry.npmmirror.com
RUN npm install -g pnpm
RUN pnpm install
RUN pnpm build:docker

FROM node:20-alpine

COPY --from=builder /code/bbs-go/site/.output /app/bbs-go/.output

WORKDIR /app/bbs-go

EXPOSE 3000

ENV ENV=docker
CMD ["node", "./.output/server/index.mjs"]