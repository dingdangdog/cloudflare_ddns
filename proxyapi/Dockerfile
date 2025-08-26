FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制源代码
COPY proxy.go .

# 构建应用
RUN go mod init proxyapi && \
    go build -o proxy_server proxy.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/proxy_server .

# 复制配置文件
COPY config_demo.json .

# 暴露端口
EXPOSE 12322

# 运行应用
CMD ["./proxy_server"]
