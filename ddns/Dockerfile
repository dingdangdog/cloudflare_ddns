# 使用官方 Go 镜像作为基础镜像
FROM golang:alpine as builder

# 设置工作目录
WORKDIR /app

# 将源代码复制到容器内
COPY . .

# 构建 Go 应用
RUN go build -o  ddns_client ddns.go

# 使用官方 Alpine 镜像作为运行时镜像
FROM alpine:latest

WORKDIR /app

# 安装 libc 依赖
RUN apk add --no-cache libc6-compat

# 复制构建好的二进制文件到容器
COPY --from=builder /app/config_demo.json ./config.json
COPY --from=builder /app/ddns_client ./
RUN chmod +x /app/ddns_client

VOLUME /app/config.json

EXPOSE 12321
# 设定容器启动命令
CMD ["/app/ddns_client"]
