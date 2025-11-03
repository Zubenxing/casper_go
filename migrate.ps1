#!/usr/bin/env pwsh
# 数据库迁移脚本 (PowerShell)

param(
    [Parameter(Position=0)]
    [ValidateSet("up", "down", "status")]
    [string]$Action = "up",
    
    [string]$ConfigPath = "config"
)

Write-Host "=== Casper 数据库迁移工具 ===" -ForegroundColor Cyan
Write-Host "操作: $Action" -ForegroundColor Yellow

# 执行迁移
go run cmd/migrate/main.go -action=$Action -config=$ConfigPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ 操作成功完成" -ForegroundColor Green
} else {
    Write-Host "✗ 操作失败" -ForegroundColor Red
    exit $LASTEXITCODE
}

