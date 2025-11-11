# n8n 工作流集成指南

## 📖 简介

本项目集成了 [n8n](https://n8n.io) 工作流自动化平台，让你可以：
- 创建和管理自动化工作流
- 通过前端界面触发工作流执行
- 查看工作流执行结果和日志
- 集成 400+ 第三方服务
- 使用 AI 功能构建智能工作流

## 🚀 快速开始

### 1. 启动 n8n 服务

```bash
# 进入项目目录
cd casper_go

# 启动 n8n 和 PostgreSQL
docker-compose -f docker-compose.n8n.yml up -d

# 查看服务状态
docker-compose -f docker-compose.n8n.yml ps

# 查看日志
docker-compose -f docker-compose.n8n.yml logs -f n8n
```

### 2. 访问 n8n UI

打开浏览器访问：http://localhost:5678

首次访问会要求创建管理员账户。

### 3. 生成 API Key

1. 登录 n8n UI
2. 点击右上角用户头像 → Settings
3. 选择 "API" 标签
4. 点击 "Create API key"
5. 复制生成的 API key
6. 将 API key 添加到 `casper_go/config/site_info.yaml`：

```yaml
n8n:
  api_url: "http://localhost:5678"
  api_key: "你的_API_KEY"
```

## 📝 创建工作流示例

### 示例：获取 GitHub 热门项目

1. 在 n8n UI 中创建新工作流
2. 添加以下节点：

#### HTTP Request 节点
- **Method**: GET
- **URL**: `https://api.github.com/search/repositories`
- **Query Parameters**:
  - `q`: `stars:>1000`
  - `sort`: `stars`
  - `order`: `desc`
  - `per_page`: `10`

#### Code 节点 (处理数据)
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

3. 保存工作流，记下工作流 ID

### 在前端触发工作流

工作流创建后，可以通过前端"AI 工作流"模块触发执行，查看结果。

## 🔧 配置说明

### Docker Compose 配置

- **n8n UI**: http://localhost:5678
- **PostgreSQL**: localhost:5433
- **数据持久化**: Docker volumes

### 环境变量

在 `docker-compose.n8n.yml` 中配置：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| DB_TYPE | 数据库类型 | postgresdb |
| N8N_PORT | n8n 端口 | 5678 |
| GENERIC_TIMEZONE | 时区 | Asia/Shanghai |
| EXECUTIONS_DATA_SAVE_ON_SUCCESS | 保存成功执行 | all |
| N8N_API_KEY_AUTH_ENABLED | 启用 API Key 认证 | true |

## 📡 API 使用

### 通过后端 API 调用 n8n

后端提供了 n8n 的代理接口：

#### 1. 获取工作流列表

```http
GET /api/v1/workflows
Authorization: Bearer <your_jwt_token>
```

#### 2. 执行工作流

```http
POST /api/v1/workflows/:id/execute
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "data": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

#### 3. 获取执行结果

```http
GET /api/v1/workflows/executions/:executionId
Authorization: Bearer <your_jwt_token>
```

## 🎨 前端集成

前端"AI 工作流"模块提供：

1. **工作流列表**：显示所有可用的工作流
2. **触发执行**：一键执行工作流
3. **执行历史**：查看历史执行记录
4. **结果展示**：格式化显示执行结果
5. **日志查看**：查看详细执行日志

## 🔒 安全建议

1. **修改默认密码**：修改 PostgreSQL 默认密码
2. **API Key 保护**：不要将 API Key 提交到版本控制
3. **网络隔离**：生产环境使用内网访问
4. **HTTPS**：生产环境启用 HTTPS
5. **权限控制**：限制工作流执行权限

## 🐛 故障排除

### n8n 无法启动

```bash
# 查看日志
docker-compose -f docker-compose.n8n.yml logs n8n

# 重启服务
docker-compose -f docker-compose.n8n.yml restart n8n
```

### 数据库连接失败

检查 PostgreSQL 是否正常运行：

```bash
docker-compose -f docker-compose.n8n.yml ps n8n-postgres
docker-compose -f docker-compose.n8n.yml logs n8n-postgres
```

### API 调用失败

1. 确认 API Key 已正确配置
2. 检查 n8n 服务是否正常运行
3. 查看后端日志

## 📚 更多资源

- [n8n 官方文档](https://docs.n8n.io)
- [n8n API 文档](https://docs.n8n.io/api/)
- [工作流模板库](https://n8n.io/workflows)
- [n8n 社区论坛](https://community.n8n.io)

## 🔄 数据备份

### 备份 n8n 数据

```bash
# 备份 PostgreSQL 数据
docker exec casper_n8n_postgres pg_dump -U n8n_user n8n > n8n_backup.sql

# 备份 n8n 配置和加密密钥
docker cp casper_n8n:/home/node/.n8n ./n8n_config_backup
```

### 恢复数据

```bash
# 恢复数据库
cat n8n_backup.sql | docker exec -i casper_n8n_postgres psql -U n8n_user -d n8n

# 恢复配置
docker cp ./n8n_config_backup/. casper_n8n:/home/node/.n8n
```

## 📊 监控和维护

### 查看资源使用

```bash
# 查看容器资源使用
docker stats casper_n8n casper_n8n_postgres
```

### 清理旧数据

定期清理旧的执行记录：
1. 在 n8n UI 中设置执行数据保留策略
2. Settings → Execution Data → Configure retention

## 🆕 更新 n8n

```bash
# 拉取最新镜像
docker-compose -f docker-compose.n8n.yml pull

# 重启服务
docker-compose -f docker-compose.n8n.yml down
docker-compose -f docker-compose.n8n.yml up -d
```

注意：更新前请先备份数据！

