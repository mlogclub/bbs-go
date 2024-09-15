# server builder

FROM golang:1.23 AS server_builder

ENV APP_HOME=/code/bbs-go/server
WORKDIR "$APP_HOME"

# COPY ./server/go.mod ./
# COPY ./server/go.sum ./
COPY ./server ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
RUN go mod download

# COPY ./server ./
RUN CGO_ENABLED=0 go build -v -o bbs-go main.go && chmod +x bbs-go

# site builder
FROM node:20-alpine AS site_builder

ENV APP_HOME=/code/bbs-go/site
WORKDIR "$APP_HOME"

COPY ./site ./
RUN npm install -g pnpm --registry=https://registry.npmmirror.com
RUN pnpm install --registry=https://registry.npmmirror.com
RUN npm install -g pnpm
RUN pnpm install
RUN pnpm build:docker



# run
FROM node:20-alpine

ENV APP_HOME=/app/bbs-go
WORKDIR "$APP_HOME"

COPY --from=site_builder /code/bbs-go/site/.output ./site/.output
COPY --from=server_builder /code/bbs-go/server/bbs-go ./server
COPY start.sh ${APP_HOME}/start.sh


EXPOSE 8082
EXPOSE 3000

ENV ENV=docker

# CMD ["${APP_HOME}/server/bbs-go"]
# CMD ["node", "${APP_HOME}/site/.output/server/index.mjs"]

ENTRYPOINT ["${APP_HOME}/start.sh"]

