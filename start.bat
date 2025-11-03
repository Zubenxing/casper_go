@echo off
chcp 65001 >nul
echo === Casper Platform 启动脚本 ===
echo.

REM 检查配置文件
if not exist "config\database.yaml" (
    echo 错误: 配置文件 config\database.yaml 不存在
    pause
    exit /b 1
)

if not exist "config\site_info.yaml" (
    echo 错误: 配置文件 config\site_info.yaml 不存在
    pause
    exit /b 1
)

REM 检查 Go 环境
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 错误: 未找到 Go 环境，请先安装 Go
    pause
    exit /b 1
)

echo 正在检查依赖...
go mod download

echo.
echo 正在启动服务器...
echo.

REM 运行程序
go run cmd\server\main.go

pause

