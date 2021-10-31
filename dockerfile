FROM golang:alpine AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io,direct
# 移动到工作目录：/build
WORKDIR /build

COPY . .

RUN go build -o app .

EXPOSE 8080
# 需要运行的命令
ENTRYPOINT  ["/build/app"]