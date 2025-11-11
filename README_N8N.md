# 🤖 n8n 工作流平台使用指南

## 📖 简介

本项目集成了 [n8n](https://n8n.io) 工作流自动化平台，使用 **官方 Docker 镜像 + PostgreSQL 数据库**。

**功能特点：**
- 🔄 数据同步和处理
- 📧 邮件自动化
- 🤖 AI 助手和聊天机器人
- 📊 数据采集和分析
- 🔔 通知和提醒
- 🌐 API 集成和编排

---

## ⚡ 快速开始

### 启动 n8n

**一键启动（推荐）：**
```powershell
.\n8n-start.ps1
```

**或使用 docker-compose：**
```bash
docker-compose -f docker-compose.n8n.yml up -d
```

### 停止 n8n

```powershell
.\n8n-stop.ps1
```

或

```bash
docker-compose -f docker-compose.n8n.yml down
```

**访问地址：** http://localhost:5678

---

## 🔧 初始化配置

### 1. 创建管理员账户

1. 打开浏览器访问：http://localhost:5678
2. 首次访问会要求创建管理员账户
3. 填写邮箱和密码完成注册

### 2. 生成 API Key（用于后端集成）

1. 登录 n8n 后，点击右上角头像
2. 选择 **Settings**
3. 点击左侧的 **API** 标签
4. 点击 **Create API key** 按钮
5. 复制生成的 API key
6. 打开 `config/site_info.yaml`，找到 n8n 配置部分：

```yaml
n8n:
  api_url: "http://localhost:5678"
  api_key: "n8n_api_YOUR_GENERATED_KEY_HERE"  # 粘贴你的 API key
```

7. 重启后端服务使配置生效

---

## 📦 服务架构

```
┌─────────────────┐
│   n8n UI/API    │  端口: 5678
│  (官方镜像)     │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│   PostgreSQL    │  端口: 5433
│    数据库       │
└─────────────────┘
```

**配置说明：**
- n8n 使用官方镜像 `n8nio/n8n:latest`
- PostgreSQL 16 Alpine（轻量级）
- 数据持久化到 Docker volumes
- 时区：Asia/Shanghai

## 📝 创建第一个工作流

### 示例：GitHub 热门项目收集器

1. 登录 n8n (http://localhost:5678)
2. 点击 **Add workflow** 创建新工作流
3. 点击 **+** 添加节点

#### 步骤 1: HTTP Request 节点

1. 搜索并添加 **HTTP Request** 节点
2. 配置如下：
   - **Method**: GET
   - **URL**: `https://api.github.com/search/repositories`
   - **Query Parameters**:
     - `q`: `stars:>5000`
     - `sort`: `stars`
     - `order`: `desc`
     - `per_page`: `10`

#### 步骤 2: Code 节点（处理数据）

1. 添加 **Code** 节点
2. 粘贴以下代码：

```javascript
// 提取需要的字段
const items = $input.all()[0].json.items;

return items.map(item => ({
  json: {
    name: item.name,
    full_name: item.full_name,
    description: item.description,
    stars: item.stargazers_count,
    forks: item.forks_count,
    language: item.language,
    url: item.html_url,
    created_at: item.created_at,
    updated_at: item.updated_at
  }
}));
```

3. 点击右上角 **Save** 保存工作流
4. 输入工作流名称：`GitHub 热门项目`
5. 点击 **Execute Workflow** 测试运行

## 🎯 在前端使用工作流

### 查看工作流列表

1. 登录系统
2. 点击左侧菜单的 **AI 工作流**
3. 可以看到所有创建的工作流

### 执行工作流

1. 在工作流卡片上点击 **执行** 按钮
2. 可选：输入 JSON 格式的参数
3. 点击 **执行** 开始运行
4. 查看执行结果和输出数据

### 查看执行历史

在 **执行历史** 表格中可以查看：
- 工作流名称
- 执行状态（成功/失败/运行中）
- 开始和结束时间
- 详细的执行结果

## 🔧 常用工作流模板

### 1. 每日天气提醒

```yaml
触发器: Schedule (每天早上 8:00)
↓
HTTP Request: 获取天气 API
↓
Code: 格式化天气信息
↓
Email/钉钉/企业微信: 发送通知
```

### 2. RSS 订阅聚合

```yaml
触发器: Schedule (每小时)
↓
RSS Feed Read: 读取多个 RSS 源
↓
Code: 去重和排序
↓
数据库/文件: 保存新文章
```

### 3. 数据库备份

```yaml
触发器: Schedule (每天凌晨 2:00)
↓
Execute Command: 执行备份命令
↓
FTP/云存储: 上传备份文件
↓
Email: 发送备份报告
```

### 4. AI 内容摘要

```yaml
触发器: Webhook
↓
HTTP Request: 获取文章内容
↓
OpenAI/Claude: 生成摘要
↓
数据库: 保存摘要
↓
返回结果
```

## 🔌 API 调用示例

### 后端 API 端点

```
GET  /api/workflows              - 获取工作流列表
GET  /api/workflows/:id          - 获取工作流详情
POST /api/workflows/:id/execute  - 执行工作流
GET  /api/workflows/executions   - 获取执行历史
GET  /api/workflows/executions/:id - 获取执行详情
POST /api/workflows/:id/activate - 激活工作流
POST /api/workflows/:id/deactivate - 停用工作流
```

### 使用 curl 执行工作流

```bash
# 获取 JWT token
TOKEN="your_jwt_token"

# 执行工作流
curl -X POST http://localhost:8080/api/workflows/WORKFLOW_ID/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"param1": "value1"}'
```

### 使用 JavaScript 调用

```javascript
import { executeWorkflow } from '@/features/workflows/api'

// 执行工作流
const result = await executeWorkflow('workflow_id', {
  query: 'stars:>1000',
  language: 'go'
})

console.log(result.data)
```

---

## 📋 常用命令

### 服务管理
```bash
# 启动服务
.\n8n-start.ps1
# 或
docker-compose -f docker-compose.n8n.yml up -d

# 停止服务
.\n8n-stop.ps1
# 或
docker-compose -f docker-compose.n8n.yml down

# 重启服务
docker-compose -f docker-compose.n8n.yml restart

# 查看服务状态
docker-compose -f docker-compose.n8n.yml ps
```

### 日志查看
```bash
# 查看所有服务日志
docker-compose -f docker-compose.n8n.yml logs -f

# 只看 n8n 日志
docker-compose -f docker-compose.n8n.yml logs -f n8n

# 只看数据库日志
docker-compose -f docker-compose.n8n.yml logs -f n8n-postgres
```

### 数据备份
```bash
# 备份 PostgreSQL 数据库
docker exec casper_n8n_postgres pg_dump -U n8n_user n8n > n8n_backup.sql

# 备份 n8n 配置
docker cp casper_n8n:/home/node/.n8n ./n8n_config_backup

# 恢复数据库
docker exec -i casper_n8n_postgres psql -U n8n_user -d n8n < n8n_backup.sql
```

### 更新镜像
```bash
# 拉取最新 n8n 镜像
docker-compose -f docker-compose.n8n.yml pull n8n

# 重新创建容器
docker-compose -f docker-compose.n8n.yml up -d --force-recreate n8n
```

---

## 🛠️ 常见问题

### Q1: n8n 无法启动？
检查以下几点：
- Docker 服务是否运行
- 端口 5678 和 5433 是否被占用
- 查看日志：`docker-compose -f docker-compose.n8n.yml logs`

### Q2: 如何修改数据库密码？
编辑 `docker-compose.n8n.yml`，修改 `POSTGRES_PASSWORD` 和 `DB_POSTGRESDB_PASSWORD`，然后重新创建容器。

### Q3: 如何访问数据库？
```bash
docker exec -it casper_n8n_postgres psql -U n8n_user -d n8n
```

### Q4: 端口冲突怎么办？
修改 `docker-compose.n8n.yml` 中的端口映射，例如将 `5678:5678` 改为 `5679:5678`。

---

## 📚 参考资源

- **[n8n 官方文档](https://docs.n8n.io)** - 完整的使用指南
- **[工作流模板库](https://n8n.io/workflows)** - 现成的工作流模板
- **[n8n 节点文档](https://docs.n8n.io/integrations/builtin/app-nodes/)** - 所有节点的使用说明

---

## 💡 最佳实践

- ✅ 使用清晰的工作流和节点名称
- ✅ 添加错误处理和重试机制
- ✅ 定期导出工作流备份
- ✅ 修改默认数据库密码
- ✅ 生产环境启用 HTTPS
- ✅ 不要将 API Key 提交到代码仓库

---

**Happy Automation! 🚀**

