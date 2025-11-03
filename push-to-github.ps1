#!/usr/bin/env pwsh
# 快速推送到 GitHub 脚本

param(
    [string]$CommitMessage = "Update",
    [switch]$Init = $false
)

Write-Host "=== Casper GitHub 推送工具 ===" -ForegroundColor Cyan
Write-Host ""

# 检查是否安装了 Git
if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "❌ 错误: 未检测到 Git" -ForegroundColor Red
    Write-Host ""
    Write-Host "请先安装 Git:" -ForegroundColor Yellow
    Write-Host "  下载地址: https://git-scm.com/download/win" -ForegroundColor Yellow
    Write-Host ""
    exit 1
}

Write-Host "✓ 检测到 Git: $(git --version)" -ForegroundColor Green

# 如果是首次初始化
if ($Init) {
    Write-Host ""
    Write-Host "=== 首次初始化 ===" -ForegroundColor Yellow
    Write-Host ""
    
    # 检查是否已经初始化
    if (Test-Path ".git") {
        Write-Host "⚠ Git 仓库已存在，跳过初始化" -ForegroundColor Yellow
    } else {
        Write-Host "初始化 Git 仓库..." -ForegroundColor Cyan
        git init
        Write-Host "✓ Git 仓库初始化完成" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "请按以下步骤操作:" -ForegroundColor Yellow
    Write-Host "1. 在 GitHub 创建新仓库: https://github.com/new" -ForegroundColor White
    Write-Host "2. 复制仓库 URL (例如: https://github.com/username/casper_go.git)" -ForegroundColor White
    Write-Host ""
    $repoUrl = Read-Host "请输入 GitHub 仓库 URL"
    
    if ($repoUrl) {
        Write-Host ""
        Write-Host "添加远程仓库..." -ForegroundColor Cyan
        git remote add origin $repoUrl
        Write-Host "✓ 远程仓库添加成功" -ForegroundColor Green
        Write-Host ""
        Write-Host "远程仓库: $repoUrl" -ForegroundColor Green
    } else {
        Write-Host "❌ 未输入仓库 URL，请手动添加:" -ForegroundColor Red
        Write-Host "  git remote add origin <your-repo-url>" -ForegroundColor Yellow
        exit 1
    }
}

# 检查是否有远程仓库
$remotes = git remote
if (-not $remotes) {
    Write-Host ""
    Write-Host "❌ 错误: 未配置远程仓库" -ForegroundColor Red
    Write-Host ""
    Write-Host "请运行以下命令添加远程仓库:" -ForegroundColor Yellow
    Write-Host "  git remote add origin https://github.com/username/casper_go.git" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "或使用初始化模式:" -ForegroundColor Yellow
    Write-Host "  .\push-to-github.ps1 -Init" -ForegroundColor Yellow
    Write-Host ""
    exit 1
}

Write-Host ""
Write-Host "=== 开始推送 ===" -ForegroundColor Yellow
Write-Host ""

# 查看状态
Write-Host "检查文件变更..." -ForegroundColor Cyan
$status = git status --short
if (-not $status) {
    Write-Host "✓ 没有需要提交的变更" -ForegroundColor Green
    Write-Host ""
    exit 0
}

Write-Host "发现以下变更:" -ForegroundColor Yellow
git status --short
Write-Host ""

# 添加所有文件
Write-Host "添加文件到暂存区..." -ForegroundColor Cyan
git add .
Write-Host "✓ 文件已添加" -ForegroundColor Green
Write-Host ""

# 提交
Write-Host "提交变更..." -ForegroundColor Cyan
Write-Host "提交消息: $CommitMessage" -ForegroundColor White
git commit -m $CommitMessage
Write-Host "✓ 提交完成" -ForegroundColor Green
Write-Host ""

# 推送
Write-Host "推送到 GitHub..." -ForegroundColor Cyan
$branch = git branch --show-current
if (-not $branch) {
    $branch = "main"
}

Write-Host "当前分支: $branch" -ForegroundColor White

try {
    git push -u origin $branch 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✓ 推送成功!" -ForegroundColor Green
        Write-Host ""
        
        # 获取远程 URL
        $remoteUrl = git remote get-url origin
        if ($remoteUrl -match "github.com[:/](.+?)\.git") {
            $repoPath = $matches[1]
            Write-Host "查看你的项目: https://github.com/$repoPath" -ForegroundColor Cyan
        }
        Write-Host ""
    } else {
        Write-Host ""
        Write-Host "❌ 推送失败" -ForegroundColor Red
        Write-Host ""
        Write-Host "可能的原因:" -ForegroundColor Yellow
        Write-Host "1. 需要登录 - 请输入 GitHub 用户名和 Token (不是密码!)" -ForegroundColor White
        Write-Host "2. 远程有新提交 - 先执行: git pull origin $branch" -ForegroundColor White
        Write-Host "3. 没有权限 - 检查 Token 是否有 repo 权限" -ForegroundColor White
        Write-Host ""
        Write-Host "如何创建 Token:" -ForegroundColor Yellow
        Write-Host "  1. 访问: https://github.com/settings/tokens" -ForegroundColor White
        Write-Host "  2. Generate new token (classic)" -ForegroundColor White
        Write-Host "  3. 勾选 'repo' 权限" -ForegroundColor White
        Write-Host "  4. 生成并复制 Token" -ForegroundColor White
        Write-Host "  5. 推送时粘贴 Token 作为密码" -ForegroundColor White
        Write-Host ""
        exit 1
    }
} catch {
    Write-Host ""
    Write-Host "❌ 推送时发生错误: $_" -ForegroundColor Red
    Write-Host ""
    exit 1
}

