# 🔍 n8n "not found" 错误调试指南

## ❌ 问题

执行工作流时返回 500 错误：`执行工作流失败: not found`

---

## 🎯 可能的原因

### 1. 工作流未正确保存

**检查步骤：**
1. 打开 n8n：http://localhost:5678
2. 找到你的工作流 "My workflow"
3. 确认是否已经点击 **Save** 按钮保存
4. 检查工作流是否在工作流列表中

### 2. 工作流 ID 不匹配

**检查步骤：**
1. 在 n8n 中打开你的工作流
2. 查看浏览器地址栏，URL 格式：
   ```
   http://localhost:5678/workflow/{workflow-id}
   ```
3. 记下这个 `workflow-id`
4. 在前端控制台查看执行时使用的 ID 是否一致

### 3. n8n API 未正确配置

**检查步骤：**
1. 检查 `config/site_info.yaml`：
   ```yaml
   n8n:
     api_url: "http://localhost:5678"
     api_key: "n8n_api_xxxxxxxxxx"  # 确保已配置
   ```
2. 确认 API Key 已经生成并正确填写

### 4. n8n API 执行端点问题

n8n 的执行工作流 API 有两种方式：

**方式 A：通过 Webhook 触发**
- 工作流需要有 Webhook 节点作为触发器
- URL：`http://localhost:5678/webhook/{webhook-path}`

**方式 B：通过 API 直接执行**
- 需要工作流有"Manual"触发器或可以手动执行
- URL：`http://localhost:5678/api/v1/workflows/{id}/execute`

---

## 🔧 解决方案

### 方案 1：确保工作流可执行

1. **在 n8n 中打开工作流**
2. **检查触发节点**：
   - 你的工作流应该有一个 **"When clicking 'Execute workflow'"** 触发器
   - 或者添加一个 **Manual Trigger** 节点

3. **测试手动执行**：
   - 在 n8n 编辑器中点击右下角的 **"Test workflow"** 按钮
   - 如果能成功执行，说明工作流本身没问题

### 方案 2：使用 Webhook 触发（推荐）

如果 API 执行有问题，可以改用 Webhook：

1. **在 n8n 中：**
   - 删除或禁用现有的触发节点
   - 添加 **Webhook** 节点
   - 配置：
     - HTTP Method: `POST`
     - Path: `my-workflow` (自定义路径)
   - 保存工作流

2. **修改前端 API（需要修改后端）：**
   - URL 从 `/workflows/{id}/execute` 改为调用 webhook
   - 地址：`http://localhost:5678/webhook/my-workflow`

### 方案 3：检查 n8n API Key 权限

1. 在 n8n 中重新生成 API Key：
   - Settings → API → Create new API key
   - 确保 API Key 有执行工作流的权限

2. 更新 `config/site_info.yaml`

3. 重启后端服务

---

## 🧪 快速测试

### 测试 1：检查 n8n API 连接

在浏览器或 Postman 中测试：

```bash
# 获取工作流列表
GET http://localhost:5678/api/v1/workflows
Headers:
  X-N8N-API-KEY: your_api_key_here
```

如果返回工作流列表，说明 API Key 有效。

### 测试 2：手动执行工作流

```bash
# 执行工作流
POST http://localhost:5678/api/v1/workflows/{workflow-id}/execute
Headers:
  X-N8N-API-KEY: your_api_key_here
Content-Type: application/json
Body: {}
```

**预期结果：**
- 成功：返回 200 和执行结果
- 失败：返回错误信息

---

## 💡 调试技巧

### 1. 查看前端控制台

刷新页面后再次执行工作流，查看控制台输出：
```
准备执行工作流: {id: "xxx", name: "My workflow", params: {}}
```

记下这个 `id`，然后在 n8n 中确认这个 ID 是否存在。

### 2. 查看后端日志

如果使用 `air` 或直接运行 Go 服务，查看终端输出：
```
[INFO] 工作流 xxx 执行成功
[ERROR] 执行工作流失败: not found
```

### 3. 检查 n8n 日志

如果使用 Docker 运行 n8n：
```bash
docker logs casper_n8n
```

查看是否有相关错误信息。

---

## 🎯 最可能的原因

根据错误信息 `not found`，最可能的原因是：

**工作流的执行端点不正确**

n8n 的工作流有几种触发方式：
1. **Manual** - 手动执行（在编辑器中点击测试）
2. **Webhook** - 通过 HTTP 请求触发
3. **Schedule** - 定时执行
4. **External** - 外部服务触发

如果你的工作流使用的是 "When clicking 'Execute workflow'" 触发器，这个触发器**只能在 n8n 编辑器中使用**，不能通过 API 调用。

---

## ✅ 推荐解决方案

**改用 Webhook 触发器：**

### 步骤 1：修改 n8n 工作流

1. 打开 n8n：http://localhost:5678
2. 打开 "My workflow"
3. **删除** "When clicking 'Execute workflow'" 节点
4. **添加** "Webhook" 节点：
   - 点击 `+` → Search "Webhook"
   - 配置：
     - **HTTP Method**: `POST` 或 `GET`
     - **Path**: `github-trending` (自定义)
   - 其他保持默认
5. 连接 Webhook 到后续节点
6. 点击 **Save**
7. 复制 Webhook URL（形如：`http://localhost:5678/webhook/github-trending`）

### 步骤 2：前端直接调用 Webhook

在前端可以直接调用这个 Webhook URL：

```javascript
// 不需要修改代码，只需知道 Webhook URL
// 或者可以存储 Webhook URL 到工作流配置中
```

### 步骤 3：测试

1. 在浏览器或 Postman 中访问：
   ```
   http://localhost:5678/webhook/github-trending
   ```
2. 应该能看到工作流执行结果
3. 如果成功，前端就可以直接调用这个 URL

---

## 📞 下一步

1. **先在 n8n 中测试** 工作流是否能手动执行
2. **查看控制台输出** 的工作流 ID
3. **确认 API Key** 已正确配置
4. **考虑使用 Webhook** 替代 API 执行

告诉我测试结果，我会帮你进一步解决！
