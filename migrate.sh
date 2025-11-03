#!/bin/bash
# 数据库迁移脚本 (Bash)

ACTION=${1:-up}
CONFIG_PATH=${2:-config}

echo "=== Casper 数据库迁移工具 ==="
echo "操作: $ACTION"

# 执行迁移
go run cmd/migrate/main.go -action="$ACTION" -config="$CONFIG_PATH"

if [ $? -eq 0 ]; then
    echo "✓ 操作成功完成"
else
    echo "✗ 操作失败"
    exit 1
fi

