FROM node:20 AS builder

ENV APP_HOME /code/bbs-go/site
WORKDIR "$APP_HOME"

COPY . ./
RUN pnpm i \
    pnpm build

FROM node:20

COPY --from=builder /code/bbs-go/site/bbs-go /app/bbs-go

WORKDIR /app

EXPOSE 8082

# 设置环境变量
ENV ENV=docker
CMD ["./bbs-go"]