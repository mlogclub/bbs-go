FROM golang:1.22 AS builder

ENV APP_HOME /code/bbs-go/server
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
# RUN go env -w GOPROXY=https://goproxy.cn,direct \
#     && go mod download
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -v -o bbs-go main.go && chmod +x bbs-go

FROM alpine:latest

COPY --from=builder /code/bbs-go/server/bbs-go /app/bbs-go
COPY --from=builder /code/bbs-go/server/bbs-go.docker.yaml /app/bbs-go.docker.yaml

WORKDIR /app

EXPOSE 8082
ENV ENV=docker
CMD ["./bbs-go"]