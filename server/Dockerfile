FROM golang:alpine AS builder

ENV APP_HOME /code/bbs-go/server
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

COPY . ./
RUN go build

EXPOSE 8082
CMD ["./bbs-go", "--config", "bbs-go.docker.yaml"]