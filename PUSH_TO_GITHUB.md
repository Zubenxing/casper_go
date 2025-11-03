# 推送项目到 GitHub 指南

## 前置准备

### 1. 安装 Git

如果还没有安装 Git，请先安装：

#### Windows
1. 下载：https://git-scm.com/download/win
2. 运行安装程序，一路默认即可
3. 安装完成后重启 PowerShell/CMD

#### 验证安装
```bash
git --version
# 应该输出类似：git version 2.x.x
```

### 2. 配置 Git

首次使用需要配置用户信息：

```bash
git config --global user.name "你的GitHub用户名"
git config --global user.email "你的GitHub邮箱"
```

### 3. 创建 GitHub 账号

如果还没有 GitHub 账号：
1. 访问 https://github.com
2. 点击 "Sign up" 注册
3. 验证邮箱

## 推送步骤

### 步骤 1: 在 GitHub 创建新仓库

1. 登录 GitHub
2. 点击右上角的 `+` → `New repository`
3. 填写信息：
   - **Repository name**: `casper_go`（或你喜欢的名字）
   - **Description**: `SSL证书监控和管理平台`
   - **Visibility**: 选择 `Public`（公开）或 `Private`（私有）
   - ⚠️ **不要勾选** "Initialize this repository with a README"
4. 点击 `Create repository`

### 步骤 2: 初始化本地仓库

在项目目录打开 PowerShell：

```powershell
# 进入项目目录
cd E:\cursor-play\casper_go

# 初始化 Git 仓库
git init

# 查看状态
git status
```

### 步骤 3: 添加文件到暂存区

```bash
# 添加所有文件（.gitignore 会自动过滤不需要的文件）
git add .

# 查看将要提交的文件
git status
```

### 步骤 4: 提交到本地仓库

```bash
# 首次提交
git commit -m "Initial commit: SSL证书监控平台 v1.0.0

- ✅ SSL证书监控功能
- ✅ 证书工具（CSR生成/验证）
- ✅ 密码管理功能
- ✅ 数据库迁移系统
- ✅ 并发优化（性能提升8+倍）
- ✅ 前端体验优化"
```

### 步骤 5: 关联远程仓库

替换 `your-username` 为你的 GitHub 用户名：

```bash
git remote add origin https://github.com/your-username/casper_go.git
```

### 步骤 6: 推送到 GitHub

```bash
# 推送到主分支
git push -u origin main

# 或者如果默认分支是 master
git push -u origin master
```

如果提示需要登录：
- 用户名：你的 GitHub 用户名
- 密码：使用 **Personal Access Token**（见下方说明）

### 步骤 7: 刷新 GitHub 页面

打开你的仓库页面，应该能看到所有代码了！

## 使用 Personal Access Token (PAT)

GitHub 已不再支持密码登录，需要使用 Token：

### 创建 Token

1. 登录 GitHub
2. 点击右上角头像 → `Settings`
3. 左侧菜单 → `Developer settings`
4. `Personal access tokens` → `Tokens (classic)`
5. 点击 `Generate new token (classic)`
6. 填写信息：
   - **Note**: `casper_go`
   - **Expiration**: 选择过期时间
   - **Select scopes**: 勾选 `repo`（所有子选项）
7. 点击 `Generate token`
8. ⚠️ **复制并保存 Token**（只显示一次！）

### 使用 Token

当 `git push` 要求输入密码时：
- 用户名：你的 GitHub 用户名
- 密码：**粘贴刚才复制的 Token**

## 常见问题

### 问题 1: 推送被拒绝 (rejected)

```bash
# 解决方法：先拉取远程代码
git pull origin main --allow-unrelated-histories
git push origin main
```

### 问题 2: 默认分支名称不匹配

如果远程是 `main`，本地是 `master`（或相反）：

```bash
# 查看当前分支
git branch

# 重命名分支
git branch -M main

# 推送
git push -u origin main
```

### 问题 3: 文件太大无法推送

如果有大文件（>100MB）：

```bash
# 查看大文件
find . -size +50M

# 添加到 .gitignore
echo "path/to/large/file" >> .gitignore

# 从 Git 历史中移除
git rm --cached path/to/large/file
git commit --amend
```

### 问题 4: 敏感信息已提交

如果不小心提交了密码等敏感信息：

```bash
# 从历史中移除文件
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch config/database.yaml" \
  --prune-empty --tag-name-filter cat -- --all

# 强制推送
git push origin --force --all
```

⚠️ **注意**：强制推送会覆盖远程历史，谨慎使用！

## 后续更新

以后修改代码后，推送到 GitHub：

```bash
# 1. 查看修改
git status

# 2. 添加修改的文件
git add .

# 3. 提交
git commit -m "feat: 添加新功能"

# 4. 推送
git push
```

### Commit 消息规范

建议使用以下前缀：

- `feat:` 新功能
- `fix:` 修复 bug
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 代码重构
- `perf:` 性能优化
- `test:` 测试相关
- `chore:` 构建/工具配置

例如：
```bash
git commit -m "feat: 添加邮件提醒功能"
git commit -m "fix: 修复证书检查并发问题"
git commit -m "docs: 更新安装文档"
```

## 使用 SSH 密钥（推荐）

如果不想每次都输入 Token，可以配置 SSH：

### 1. 生成 SSH 密钥

```bash
ssh-keygen -t ed25519 -C "your_email@example.com"
# 一路回车，使用默认路径和空密码
```

### 2. 添加到 GitHub

```bash
# 复制公钥内容
cat ~/.ssh/id_ed25519.pub
# 或 Windows:
type %USERPROFILE%\.ssh\id_ed25519.pub
```

1. 登录 GitHub
2. `Settings` → `SSH and GPG keys`
3. 点击 `New SSH key`
4. 粘贴公钥，保存

### 3. 更改远程 URL

```bash
# 从 HTTPS 改为 SSH
git remote set-url origin git@github.com:your-username/casper_go.git

# 验证
git remote -v
```

以后推送就不需要输入密码了！

## 创建分支

开发新功能时建议使用分支：

```bash
# 创建并切换到新分支
git checkout -b feature/email-notification

# 开发完成后合并回主分支
git checkout main
git merge feature/email-notification

# 推送分支到远程
git push origin feature/email-notification
```

## 忽略已跟踪的文件

如果文件已经被 Git 跟踪，添加到 `.gitignore` 后还需要：

```bash
# 从 Git 移除（但保留本地文件）
git rm --cached config/database.yaml

# 提交
git commit -m "chore: 停止跟踪配置文件"

# 推送
git push
```

## 总结

✅ 完整流程：

1. 安装并配置 Git
2. 在 GitHub 创建仓库
3. 初始化本地仓库：`git init`
4. 添加文件：`git add .`
5. 提交：`git commit -m "Initial commit"`
6. 关联远程：`git remote add origin <url>`
7. 推送：`git push -u origin main`

🎉 **完成！你的项目现在已经在 GitHub 上了！**

分享链接：`https://github.com/your-username/casper_go`

