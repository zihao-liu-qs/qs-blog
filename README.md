# qs-blog

个人博客网站,使用 Hugo + Go + Docker 构建和部署。

## 项目结构

```
├── content/           # Hugo 内容文件
├── server/            # Go 后端服务(激活码管理)
├── themes/minimal/    # 博客主题
├── .env               # 环境配置(不提交到 Git)
├── .env.example       # 环境配置模板
├── Dockerfile         # 多阶段构建 Docker 镜像
├── docker-compose.yml # Docker Compose 配置
└── nginx.conf         # Nginx 配置
```

## 功能特性

- **博客系统**: 基于 Hugo 的静态网站生成
- **激活验证**: 软件激活码管理服务
- **Stripe 支付**: 集成 Stripe Checkout 支付
- **Docker 部署**: 一键部署,包含前端和后端

## 部署

### 1. 配置环境变量

复制配置模板并修改：

```bash
cp .env.example .env
```

编辑 `.env`，修改以下关键配置：

```env
# 基础 URL（只需改这一行，所有配置自动同步）
# 本地开发：http://localhost/
# 云服务器：http://142.93.127.127/
BASE_URL=http://localhost/

# Stripe 支付配置
STRIPE_SECRET_KEY=sk_test_xxx
STRIPE_PUBLISHABLE_KEY=pk_test_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx

# 管理接口密钥
ADMIN_KEY=change-me-in-production
```

### 2. 使用 Makefile 部署

```bash
# 部署到本地
make deploy-local

# 部署到云服务器
make deploy-cloud
```

### 3. 或使用 Docker Compose 手动部署

```bash
# 构建并启动
docker compose up -d --build

# 查看日志
docker compose logs -f

# 停止服务
docker compose down
```

访问 `http://localhost` 即可查看博客。

## 服务说明

### 前端(端口 80)
- Hugo 生成的静态网站
- Nginx 提供高性能访问

### 后端(端口 8080,内部)
- Go + Gin 框架
- SQLite 数据库
- 激活码管理和验证
- Stripe 支付集成
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
```

#### 创建支付订单
```bash
POST /api/v1/checkout
Content-Type: application/json

{
  "product_id": "blink",
  "product_name": "Blink",
  "product_price": 6800,
  "customer_email": "user@example.com",
  "success_url": "http://localhost/payment-success",
  "cancel_url": "http://localhost/software/blink"
}
```

#### Stripe Webhook 回调
```bash
POST /api/v1/webhook/stripe
```

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
| `BASE_URL` | `http://localhost/` | 网站基础 URL，Hugo 构建和运行时使用 |
| `PORT` | `8080` | 后端监听端口 |
| `ADMIN_KEY` | `change-me-in-production` | 管理接口密钥 |
| `DB_PATH` | `/data/activate.db` | 数据库路径 |
| `LOG_PATH` | `/var/log/activity.log` | 活动日志路径 |
| `LOG_MAX_AGE_DAYS` | `30` | 日志保留天数 |
| `STRIPE_SECRET_KEY` | - | Stripe Secret Key |
| `STRIPE_PUBLISHABLE_KEY` | - | Stripe Publishable Key |
| `STRIPE_WEBHOOK_SECRET` | - | Stripe Webhook 签名密钥 |

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

1. **不要提交 `.env` 文件**: 已添加到 `.gitignore`，请勿将密钥提交到 Git
2. **修改 ADMIN_KEY**: 部署前务必修改 `.env` 中的 `ADMIN_KEY`
3. **使用 HTTPS**: 生产环境建议使用反向代理配置 HTTPS
4. **定期备份**: 定期备份 `/data` 目录中的数据库文件
