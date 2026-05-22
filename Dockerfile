# syntax=docker/dockerfile:1.7

FROM node:24-alpine AS web-builder
WORKDIR /src/web

RUN corepack enable && corepack prepare pnpm@10.30.2 --activate
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml web/.npmrc ./
RUN pnpm install --frozen-lockfile

COPY web/ ./
RUN pnpm build:ssr \
	&& cp -R build build-ssr \
	&& pnpm build:spa \
	&& cp -R build/spa build-spa \
	&& rm -rf build \
	&& mv build-ssr build \
	&& mv build-spa build/spa
RUN pnpm prune --prod

FROM golang:1.26-alpine AS server-builder
WORKDIR /src

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY . ./
COPY --from=web-builder /src/web/build/spa ./web/build/spa

ARG TARGETOS=linux
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -v -trimpath -ldflags="-s -w" -o /out/bbs-go ./main.go

FROM node:24-alpine AS app
WORKDIR /app

ENV NODE_ENV=production \
	PORT=3000 \
	BBSGO_ENV=prod \
	BBSGO_SERVER_URL=http://127.0.0.1:8082 \
	TZ=Asia/Shanghai

RUN apk add --no-cache ca-certificates tzdata wget \
	&& mkdir -p /app/data /app/logs /app/res/uploads /app/defaults

COPY --from=server-builder /out/bbs-go /app/bbs-go
COPY locales /app/locales
COPY res /app/res
COPY --from=web-builder /src/web/package.json /app/package.json
COPY --from=web-builder /src/web/node_modules /app/node_modules
COPY --from=web-builder /src/web/build /app/build
COPY --from=web-builder /src/web/scripts /app/scripts
COPY docker/bbs-go-docker.yaml /app/defaults/bbs-go.yaml
COPY docker/entrypoint.sh /app/entrypoint.sh

RUN ln -s /app/data/bbs-go.yaml /app/bbs-go.yaml \
	&& chmod +x /app/entrypoint.sh

EXPOSE 3000 8082
VOLUME ["/app/data", "/app/logs", "/app/res/uploads"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
	CMD wget -qO- http://127.0.0.1:3000/api/install/status >/dev/null || exit 1

CMD ["/app/entrypoint.sh"]
