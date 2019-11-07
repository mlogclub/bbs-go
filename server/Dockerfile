# builder stage
FROM golang:1.13-stretch AS builder

# Install build dependencies
RUN apt-get -qq update && \
    apt-get -qq install -y --no-install-recommends \
      build-essential \
      git \
      openssh-client \
    && rm -rf /var/lib/apt/lists/*

# Update timezone
ENV TZ=Asia/Shanghai

WORKDIR /app

ENV GO111MODULE=on
ENV ROOT_DIR=/app

# download and cache go dependencies
COPY go.* ./
RUN GOPROXY="https://goproxy.cn" go mod download

COPY . .

RUN go build -o bbs-go-server

# application stage
FROM debian:stretch-slim as application

WORKDIR /app

# Install runtime dependencies
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl \
    && apt-get -qq install -y host mysql-client\
    && rm -rf /var/lib/apt/lists/*

# Update timezone
ENV TZ=Asia/Shanghai
ENV ROOT_DIR=/app

COPY --from=builder /app/bbs-go-server .
COPY --from=builder /app/bbs-go.yaml bbs-go.yaml

EXPOSE 8080

HEALTHCHECK --start-period=10s \
            --interval=15s \
            --timeout=5s \
            --retries=3 \
            CMD curl -sSf http://localhost:8082/api/img/proxy || exit 1

CMD ["./bbs-go-server"]
