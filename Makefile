.PHONY: deploy-local deploy-cloud

# 部署到本地 (localhost)
deploy-local:
	BASE_URL=http://localhost/ docker compose up -d --build

# 部署到云服务器
deploy-cloud:
	BASE_URL=http://<your-cloud-server-ip>/ docker compose up -d --build
