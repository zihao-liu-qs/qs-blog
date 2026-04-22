IMAGE_NAME = zihaoliuqs/qs-blog

.PHONY: deploy build-for-amd64

# 构建并部署
deploy:
	docker compose up -d --build

# 本地编译 AMD64 镜像（用于云服务器）
build-for-amd64:
	docker buildx build --platform linux/amd64 -t zihaoliuqs/qs-blog:latest --load .

# 服务器端拉取最新镜像并强制重新部署（需在服务器上执行，需存在 .env 文件）
deploy-remote:
	docker pull zihaoliuqs/qs-blog:latest
	docker stop qs-blog 2>/dev/null || true
	docker rm qs-blog 2>/dev/null || true
	docker run -d -p 80:80 --env-file .env --name qs-blog zihaoliuqs/qs-blog:latest
