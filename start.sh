#!/bin/bash

# Casper Platform 启动脚本

echo "=== Casper Platform 启动脚本 ==="
echo ""

# 检查配置文件
if [ ! -f "config/database.yaml" ]; then
    echo "错误: 配置文件 config/database.yaml 不存在"
    exit 1
fi

if [ ! -f "config/site_info.yaml" ]; then
    echo "错误: 配置文件 config/site_info.yaml 不存在"
    exit 1
fi

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go 环境，请先安装 Go"
    exit 1
fi

echo "正在检查依赖..."
go mod download

echo ""
echo "正在启动服务器..."
echo ""

# 运行程序
go run cmd/server/main.go

