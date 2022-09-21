FROM golang AS build-stage

ENV APP_HOME /code/bbs-go/server
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build

EXPOSE 8082
CMD ["./bbs-go", "--config", "bbs-go.docker.yaml"]