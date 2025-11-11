# 🚀 n8n 工作流配置指南

你的项目已经完全集成了 n8n，只需要 3 步配置即可开始使用！

## ✅ 当前状态

- ✅ n8n 服务已启动（http://localhost:5678）
- ✅ 后端 API 已实现（`features/workflow`）
- ✅ 前端页面已实现（`Workflows.vue`）
- ✅ 工作流已创建（GitHub 热门项目统计）
- ⏳ **需要配置 API Key**

---

## 🔧 配置步骤

### 步骤 1: 在 n8n 中生成 API Key

1. 打开浏览器访问：http://localhost:5678
2. 登录你的 n8n 账户
3. 点击右上角 **头像/用户名**
4. 选择 **Settings**（设置）
5. 点击左侧的 **API** 选项卡
6. 点击 **Create API key** 按钮
7. **复制** 生成的 API key（类似：`n8n_api_xxxxxxxxxxxxxxxxxxxxxxxx`）

> ⚠️ **重要**：API key 只显示一次，请立即复制保存！

---

### 步骤 2: 配置到 site_info.yaml

编辑配置文件：`config/site_info.yaml`

```yaml
# n8n 工作流引擎配置
n8n:
  api_url: "http://localhost:5678"
  api_key: "n8n_api_你刚刚复制的密钥"  # ⬅️ 粘贴你的 API key
```

**示例：**
```yaml
n8n:
  api_url: "http://localhost:5678"
  api_key: "n8n_api_1234567890abcdefghijklmnop"
```

---

### 步骤 3: 重启后端服务

配置保存后，重启 Go 后端服务：

```bash
# 如果使用 air 热重载，会自动重启
# 或手动重启：
go run cmd/server/main.go
```

---

## 🎯 使用前端触发工作流

### 方式一：通过前端 UI（推荐）

1. **访问工作流页面**
   - 打开浏览器：http://localhost:8080
   - 登录系统
   - 点击左侧菜单 **AI 工作流**

2. **查看你的工作流**
   - 你会看到 "My workflow" 工作流卡片
   - 显示工作流状态、节点数量等信息

3. **执行工作流**
   - 点击工作流卡片上的 **"执行"** 按钮
   - 弹出执行对话框
   - （可选）输入 JSON 参数
   - 点击 **"执行"** 确认

4. **查看执行结果** ⭐
   - 执行成功后，页面 **自动显示执行结果卡片**
   - 结果直接嵌入页面中（工作流列表和执行历史之间）
   - **自动滚动** 到结果区域
   - 显示执行状态、时间、返回数据
   - 可 **一键复制** 结果到剪贴板
   - 可点击 **"在 n8n 中查看详情"** 查看完整日志

5. **查看执行历史**
   - 页面下方有 **执行历史** 表格
   - 显示所有工作流的执行记录
   - 点击 **"查看"** 按钮可在页面中重新加载该执行结果
   - 可删除历史记录

---

### 方式二：通过 API 调用

如果你想在其他地方调用工作流，可以使用后端 API：

#### 1. 获取工作流 ID

**请求：**
```bash
GET http://localhost:8080/api/workflows
Authorization: Bearer YOUR_JWT_TOKEN
```

**响应：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": "workflow_id_here",
      "name": "My workflow",
      "active": true
    }
  ]
}
```

#### 2. 执行工作流

**请求：**
```bash
POST http://localhost:8080/api/workflows/{workflow_id}/execute
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "count": 20,
  "language": "all"
}
```

**响应：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "execution_id",
    "status": "success",
    "data": {
      "resultData": {
        "runData": {
          // GitHub 热门项目数据
        }
      }
    }
  }
}
```

---

## 📋 可用的后端 API

```
GET    /api/workflows                    - 获取工作流列表
GET    /api/workflows/:id                - 获取工作流详情
POST   /api/workflows/:id/execute        - 执行工作流 ⭐
GET    /api/workflows/executions         - 获取执行历史
GET    /api/workflows/executions/:id     - 获取执行详情
DELETE /api/workflows/executions/:id     - 删除执行记录
POST   /api/workflows/:id/activate       - 激活工作流
POST   /api/workflows/:id/deactivate     - 停用工作流
GET    /api/workflows/health             - 检查 n8n 连接状态
```

所有 API 都需要 JWT 认证（登录后获取 token）。

---

## 🎨 前端功能

你的前端已经实现了完整的工作流管理界面：

### 功能列表
- ✅ 显示所有工作流列表（卡片式）
- ✅ 搜索工作流
- ✅ 执行工作流（带参数输入）
- ⭐ **页面内展示执行结果**（无需弹窗，自动滚动）
- ✅ 一键复制结果到剪贴板
- ✅ 查看执行历史
- ✅ 激活/停用工作流
- ✅ 在 n8n 编辑器中打开工作流
- ✅ 检查 n8n 连接状态
- ✅ 删除执行记录

### 💡 新特性
**执行结果直接显示在页面中**，不再使用弹窗：
- 📍 结果卡片出现在工作流列表和执行历史之间
- 🎬 自动平滑滚动到结果位置
- 📋 支持一键复制 JSON 数据
- ❌ 可随时关闭结果卡片
- 🎨 美观的动画效果和格式化展示

### 页面路由
```javascript
// 前端路由：/workflows
// 组件位置：frontend/src/features/workflows/Workflows.vue
```

---

## 🔍 测试步骤

### 1. 验证配置
```bash
# 检查 n8n 连接状态
curl http://localhost:8080/api/workflows/health \
  -H "Authorization: Bearer YOUR_TOKEN"

# 应该返回：
{
  "code": 0,
  "msg": "success",
  "data": {
    "status": "healthy",
    "message": "n8n 服务运行正常"
  }
}
```

### 2. 获取工作流列表
```bash
curl http://localhost:8080/api/workflows \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 在前端执行工作流
1. 打开 http://localhost:8080/workflows
2. 找到 "My workflow"
3. 点击 "执行" 按钮
4. 查看 GitHub 热门项目数据

---

## 💡 自定义工作流参数

如果你的工作流需要参数，在执行对话框中输入 JSON：

```json
{
  "count": 20,
  "language": "go",
  "since": "weekly"
}
```

然后在 n8n 工作流中可以通过 `{{ $json.count }}` 访问这些参数。

---

## 🛠️ 常见问题

### Q1: 前端显示 "n8n 未连接"？
- 检查 n8n 服务是否运行：http://localhost:5678
- 检查 `site_info.yaml` 中的 `api_url` 是否正确
- 重启后端服务

### Q2: 执行失败 "401 Unauthorized"？
- 检查 API Key 是否正确配置
- API Key 格式：`n8n_api_xxxxxx`
- 重新生成 API Key 并更新配置

### Q3: 看不到工作流列表？
- 确保在 n8n 中已保存工作流
- 工作流名称会自动显示在列表中
- 刷新页面或点击 "刷新" 按钮

### Q4: 如何传递动态参数？
在执行对话框中输入 JSON 格式的参数，n8n 工作流可以通过 `$json.xxx` 访问。

---

## 🎉 完成！

现在你可以：
1. ✅ 在前端点击按钮执行 n8n 工作流
2. ✅ 实时查看执行结果
3. ✅ 查看所有执行历史
4. ✅ 在 n8n 中编辑工作流
5. ✅ 通过 API 在任何地方调用工作流

**就是这么简单！** 🚀

---

## 📚 相关文档

- [README_N8N.md](README_N8N.md) - n8n 基础使用指南
- [n8n 官方文档](https://docs.n8n.io)
- [n8n API 文档](https://docs.n8n.io/api/)
