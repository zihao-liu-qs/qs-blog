.PHONY: dev build

# 本地开发：热重载（需要 air：go install github.com/air-verse/air@latest）
dev:
	air -c .air.toml

# 直接运行
run:
	go run main.go

# 编译
build:
	go build -o bin/server main.go
