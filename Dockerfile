# 第一阶段：构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git make

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bin/html2md \
    cmd/server/main.go

# 第二阶段：运行时阶段
FROM alpine:3.18

# 安装运行时依赖
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    && update-ca-certificates

# 创建非root用户
RUN addgroup -g 1001 -S html2md && \
    adduser -u 1001 -S html2md -G html2md

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/html2md /app/html2md

# 创建必要的目录
RUN mkdir -p /app/logs && \
    chown -R html2md:html2md /app

# 切换到非root用户
USER html2md

# 暴露端口
EXPOSE 8080 9090

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

# 运行应用
CMD ["./html2md"] 