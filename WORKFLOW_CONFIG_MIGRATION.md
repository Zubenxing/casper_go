# 🚀 工作流配置动态化迁移指南

## 📋 概述

工作流配置已从硬编码迁移到**数据库驱动**，现在可以通过 API 动态管理，无需修改代码。

---

## ✨ 新特性

### 1. 数据库存储
- ✅ 配置存储在 `workflow_configs` 表
- ✅ 支持增删改查
- ✅ 软删除支持

### 2. API 管理
```
GET    /api/workflows/configs           # 获取所有配置
GET    /api/workflows/configs/:id       # 获取指定配置
POST   /api/workflows/configs           # 创建配置
PUT    /api/workflows/configs/:id       # 更新配置
DELETE /api/workflows/configs/:id       # 删除配置
POST   /api/workflows/configs/sync      # 从 n8n 同步
```

### 3. 自动同步
- ✅ 点击"同步"按钮自动从 n8n 拉取工作流
- ✅ 自动创建默认配置（simple 类型）
- ✅ 已存在的配置不会被覆盖

---

## 🎯 使用方式

### 方式 1：通过 API 手动添加

```javascript
// 创建配置
POST /api/workflows/configs
{
  "workflowId": "Zps0k9QUpdoXR4kD",
  "workflowName": "Azure账单对比(修正版)",
  "type": "upload_files",
  "description": "上传本月和上月的 Azure 账单 Excel 文件进行对比分析",
  "executeLabel": "开始对比",
  "resultType": "download",
  "enabled": true,
  "fields": [
    {
      "name": "current_month",
      "label": "本月账单",
      "type": "file",
      "accept": ".xlsx,.xls",
      "required": true,
      "description": "请上传本月的 Azure 账单 Excel 文件"
    },
    {
      "name": "last_month",
      "label": "上月账单",
      "type": "file",
      "accept": ".xlsx,.xls",
      "required": true,
      "description": "请上传上月的 Azure 账单 Excel 文件"
    }
  ]
}
```

### 方式 2：使用同步功能

```bash
# 从 n8n 自动同步所有工作流
POST /api/workflows/configs/sync

# 响应
{
  "synced": 3,   # 新增的配置数
  "total": 5     # 总工作流数
}
```

---

## 📦 数据库字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | uint | 主键 |
| `workflowId` | string | n8n 工作流 ID（唯一） |
| `workflowName` | string | 工作流名称 |
| `type` | string | 类型：`simple`/`upload_files`/`form_input` |
| `description` | string | 描述 |
| `executeLabel` | string | 执行按钮文本 |
| `resultType` | string | 结果类型：`auto`/`html`/`json`/`download` |
| `fields` | JSON | 字段配置数组 |
| `enabled` | bool | 是否启用 |

---

## 🔧 配置示例

### 简单执行类型
```json
{
  "workflowName": "GitHub Trending",
  "type": "simple",
  "description": "获取 GitHub 热门项目",
  "executeLabel": "获取",
  "resultType": "html",
  "fields": []
}
```

### 文件上传类型
```json
{
  "workflowName": "文件处理",
  "type": "upload_files",
  "description": "上传文件进行处理",
  "executeLabel": "开始处理",
  "resultType": "download",
  "fields": [
    {
      "name": "file",
      "label": "选择文件",
      "type": "file",
      "accept": ".pdf,.docx",
      "required": true
    }
  ]
}
```

### 表单输入类型
```json
{
  "workflowName": "数据查询",
  "type": "form_input",
  "description": "根据条件查询数据",
  "executeLabel": "查询",
  "resultType": "json",
  "fields": [
    {
      "name": "startDate",
      "label": "开始日期",
      "type": "date",
      "required": true
    },
    {
      "name": "keyword",
      "label": "关键词",
      "type": "text",
      "placeholder": "请输入关键词"
    }
  ]
}
```

---

## 🚀 启动流程

### 1. 重启后端
```bash
# 会自动创建 workflow_configs 表
go run cmd\server\main.go
```

### 2. 同步现有工作流
```bash
# 使用 Postman 或 curl
POST http://localhost:8080/api/workflows/configs/sync
Authorization: Bearer <your_token>
```

### 3. 刷新前端
前端会自动从后端加载配置

---

## 📝 添加新工作流的步骤

### 步骤 1：在 n8n 中创建工作流
1. 打开 n8n 创建工作流
2. 添加 Webhook 触发器
3. 配置节点
4. 激活工作流

### 步骤 2：同步或手动添加配置
#### 选项 A：自动同步（推荐）
```bash
POST /api/workflows/configs/sync
```

#### 选项 B：手动添加
```bash
POST /api/workflows/configs
{
  "workflowId": "xxx",
  "workflowName": "新工作流",
  ...
}
```

### 步骤 3：刷新前端
刷新浏览器，新工作流自动显示！

---

## 🎉 优势

### 之前（硬编码）
```javascript
// ❌ 每次都要改代码
export const workflowConfigs = {
  '工作流1': {...},
  '工作流2': {...}
}
```

### 现在（数据库）
```
✅ 在 n8n 创建工作流
✅ 点击"同步"按钮
✅ 完成！无需改代码
```

---

## 🔍 检查配置

```bash
# 查看所有配置
GET /api/workflows/configs

# 查看特定工作流配置
GET /api/workflows/configs/{workflowId}
```

---

## ⚡ 快速测试

```bash
# 1. 启动后端
go run cmd\server\main.go

# 2. 同步配置
curl -X POST http://localhost:8080/api/workflows/configs/sync \
  -H "Authorization: Bearer <token>"

# 3. 查看配置
curl http://localhost:8080/api/workflows/configs \
  -H "Authorization: Bearer <token>"

# 4. 刷新前端
# 打开浏览器，所有工作流配置已经动态加载！
```

---

## 🎊 完成！

现在你可以：
- ✅ 在 n8n 中自由添加/删除工作流
- ✅ 通过 API 或同步功能管理配置
- ✅ 无需修改前端代码
- ✅ 配置持久化存储在数据库

**真正实现了配置与代码分离！** 🚀
