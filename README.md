# qs-blog

个人博客网站,使用 Hugo + Go + Docker 构建和部署。

## 项目结构

```
├── content/           # Hugo 内容文件
├── server/            # Go 后端服务(激活码管理)
├── themes/minimal/    # 博客主题
├── Dockerfile         # 多阶段构建 Docker 镜像
├── docker-compose.yml # Docker Compose 配置
└── nginx.conf         # Nginx 配置
```

## 功能特性

- **博客系统**: 基于 Hugo 的静态网站生成
- **激活验证**: 软件激活码管理服务
- **Docker 部署**: 一键部署,包含前端和后端

## Docker 部署

### 使用 Docker Compose 运行(推荐)

```bash
# 构建并启动
docker compose up -d

# 查看日志
docker compose logs -f

# 停止服务
docker compose down
```

### 直接使用 Docker 运行

```bash
# 构建镜像
docker build -t qs-blog .

# 运行容器
docker run -d -p 1313:80 --name qs-blog qs-blog

# 查看日志
docker logs -f qs-blog

# 停止容器
docker stop qs-blog
```

访问 http://localhost:1313 即可查看博客。

## 服务说明

### 前端(端口 80)
- Hugo 生成的静态网站
- Nginx 提供高性能访问

### 后端(端口 8080,内部)
- Go + Gin 框架
- SQLite 数据库
- 激活码管理和验证
- IP 限流和活动日志

### API 接口

#### 客户端验证
```bash
POST /api/v1/activate
Content-Type: application/json

{
  "license": "激活码",
  "device_id": "设备标识"
}
``

#### 管理接口(需要 X-Admin-Key)
```bash
# 创建激活码
POST /api/admin/licenses

# 查询所有激活码
GET /api/admin/licenses

# 查询单个激活码
GET /api/admin/licenses/:key

# 吊销激活码
DELETE /api/admin/licenses/:key

# 移除设备绑定
DELETE /api/admin/devices/:key/:device_id
```

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `8080` | 后端监听端口 |
| `ADMIN_KEY` | `change-me-in-production` | 管理接口密钥 |
| `DB_PATH` | `/data/activate.db` | 数据库路径 |
| `LOG_PATH` | `/var/log/activity.log` | 活动日志路径 |
| `LOG_MAX_AGE_DAYS` | `30` | 日志保留天数 |

## 本地开发

### 前端开发
```bash
hugo server -D
```

### 后端开发
```bash
cd server
go run main.go
```

## 安全建议

1. **修改 ADMIN_KEY**: 部署前务必修改 docker-compose.yml 中的 `ADMIN_KEY`
2. **使用 HTTPS**: 生产环境建议使用反向代理配置 HTTPS
3. **定期备份**: 定期备份 `/data` 目录中的数据库文件
