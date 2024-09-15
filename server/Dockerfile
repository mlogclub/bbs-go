FROM golang:1.23 AS builder

ENV APP_HOME /code/bbs-go/server
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -v -o bbs-go main.go && chmod +x bbs-go

FROM alpine:latest

COPY --from=builder /code/bbs-go/server/bbs-go /app/bbs-go

WORKDIR /app

EXPOSE 8082
CMD ["./bbs-go"]