# 使用官方 Golang 镜像进行构建
FROM golang:1.25 AS builder

WORKDIR /app
COPY . .

# 下载依赖并编译
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o message-dispatch main.go

# 使用更小的基础镜像
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/message-dispatch .
COPY config/ ./config/
EXPOSE 15001

# 启动服务
ENTRYPOINT ["./message-dispatch"]