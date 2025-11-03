.PHONY: help build run clean test install dev swagger

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 安装项目依赖
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: ## 生成 Swagger 文档
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

swagger-fmt: ## 格式化 Swagger 注释
	swag fmt

build: swagger ## 编译项目
	go build -o casper_platform cmd/server/main.go

build-linux: swagger ## 编译 Linux 版本
	GOOS=linux GOARCH=amd64 go build -o casper_platform_linux cmd/server/main.go

build-windows: swagger ## 编译 Windows 版本
	GOOS=windows GOARCH=amd64 go build -o casper_platform.exe cmd/server/main.go

run: swagger ## 运行项目
	go run cmd/server/main.go

dev: ## 开发模式运行（自动重载需要安装 air）
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "请先安装 air: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

test: ## 运行测试
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## 清理编译文件
	rm -f casper_platform casper_platform.exe casper_platform_linux
	rm -f coverage.out coverage.html
	rm -rf logs/*

fmt: ## 格式化代码
	go fmt ./...

vet: ## 代码静态检查
	go vet ./...

lint: ## 代码 lint 检查（需要安装 golangci-lint）
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "请先安装 golangci-lint: https://golangci-lint.run/usage/install/"; \
	fi

migrate-up: ## 执行数据库迁移（升级）
	go run cmd/migrate/main.go -action=up

migrate-down: ## 回滚最后一次数据库迁移
	go run cmd/migrate/main.go -action=down

migrate-status: ## 查看数据库迁移状态
	go run cmd/migrate/main.go -action=status

docker-build: ## 构建 Docker 镜像
	docker build -t casper_platform:latest .

docker-run: ## 运行 Docker 容器
	docker-compose up -d

docker-stop: ## 停止 Docker 容器
	docker-compose down

all: clean install swagger build ## 完整构建流程
