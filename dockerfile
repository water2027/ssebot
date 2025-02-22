# 使用 Golang 进行编译
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN export GOPROXY=https://goproxy.cn,direct && go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

COPY config.yaml .
CMD ["./main"]
