# ====================================
# 阶段 1: 使用 Hugo 构建静态站点
# ====================================
FROM hugomods/hugo:0.143.0 AS hugo-builder

ARG HUGO_BASE_URL=http://localhost/

WORKDIR /app

# 复制所有文件（排除 server 目录）
COPY . .
RUN rm -rf /app/server

# 强制重新构建站点
RUN hugo --baseURL="${HUGO_BASE_URL}" --minify --gc --forceSyncStatic

# ====================================
# 阶段 2: 构建 Go 后端服务
# ====================================
FROM golang:1.25-alpine AS go-builder

# 安装 gcc 和其他依赖（SQLite 需要 CGO）
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build

# 复制 server 目录
COPY server/ ./

# 下载依赖
RUN go mod download

# 编译 Go 服务
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /build/server .

# ====================================
# 阶段 3: 最终运行镜像
# ====================================
FROM alpine:3.19

# 安装 nginx 和 curl
RUN apk add --no-cache nginx curl

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/http.d/default.conf

# 从 Hugo 构建阶段复制静态文件
COPY --from=hugo-builder /app/public /usr/share/nginx/html

# 从 Go 构建阶段复制二进制文件
COPY --from=go-builder /build/server /usr/local/bin/server

# 创建数据和日志目录
RUN mkdir -p /data /var/log /var/log/nginx /var/lib/nginx/logs /run

# 复制启动脚本
RUN cat > /start.sh << 'EOF'
#!/bin/sh

# 启动 Go 后端服务
/usr/local/bin/server &
SERVER_PID=$!

# 等待后端服务启动
sleep 2

# 启动 nginx
nginx -g 'daemon off;' &
NGINX_PID=$!

# 等待任一进程退出
wait $SERVER_PID $NGINX_PID
EOF

RUN chmod +x /start.sh

EXPOSE 80

CMD ["/start.sh"]
