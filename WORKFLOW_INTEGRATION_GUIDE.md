# 🚀 工作流集成指南

## 📖 概述

本系统采用**配置驱动**的方式集成 n8n 工作流，无需为每个工作流单独编写代码。

## 🏗️ 架构设计

### 核心思路

1. **工作流配置中心** - 定义每个工作流的交互方式
2. **动态执行组件** - 根据配置自动渲染 UI
3. **通用 API 层** - 统一的工作流执行接口

### 工作流类型

| 类型 | 说明 | 适用场景 | 示例 |
|-----|------|---------|------|
| **SIMPLE** | 简单执行，无需参数 | 数据获取、报表生成 | GitHub Trending |
| **UPLOAD_FILES** | 上传文件 | 文件处理、数据导入 | Azure 账单对比 |
| **FORM_INPUT** | 表单输入 | 数据查询、条件筛选 | 日期范围查询 |

---

## 📝 集成新工作流步骤

### 步骤 1：在 n8n 中创建工作流

1. 打开 n8n：http://localhost:5678
2. 创建你的工作流
3. **重要**：添加 Webhook 触发器
   - HTTP Method: GET 或 POST（根据需要）
   - Path: 有意义的路径名（如 `azure-bill`）
   - Response Mode: `When Last Node Finishes`
4. **激活工作流**
5. 复制 Production URL

### 步骤 2：在前端配置文件中注册

编辑 `frontend/src/features/workflows/workflow-configs.js`：

```javascript
export const workflowConfigs = {
  // 你的工作流名称（必须与 n8n 中的名称一致）
  'Azure账单对比': {
    type: WORKFLOW_TYPES.UPLOAD_FILES, // 工作流类型
    description: '上传本月和上月的 Azure 账单 Excel 文件进行对比分析',
    executeLabel: '开始对比',
    resultType: 'download', // 结果类型：html/json/download/auto
    
    // 如果是 UPLOAD_FILES 或 FORM_INPUT 类型，需要定义 fields
    fields: [
      {
        name: 'current_month',     // 字段名（对应 n8n 接收的参数名）
        label: '本月账单',          // 显示标签
        type: 'file',              // 字段类型
        accept: '.xlsx,.xls',      // 文件类型限制
        required: true,            // 是否必填
        description: '请上传本月的 Azure 账单 Excel 文件'
      },
      {
        name: 'last_month',
        label: '上月账单',
        type: 'file',
        accept: '.xlsx,.xls',
        required: true,
        description: '请上传上月的 Azure 账单 Excel 文件'
      }
    ]
  }
}
```

### 步骤 3：刷新前端

**就这么简单！** 刷新浏览器后：
- ✅ 工作流卡片上会显示类型标签
- ✅ 点击"执行"按钮会自动显示对应的UI
- ✅ 无需编写任何额外代码

---

## 🎯 字段类型参考

### 文件上传字段

```javascript
{
  name: 'file_field',
  label: '上传文件',
  type: 'file',
  accept: '.pdf,.docx',  // 支持的文件类型
  required: true
}
```

### 文本输入字段

```javascript
{
  name: 'keyword',
  label: '关键词',
  type: 'text',
  required: false,
  placeholder: '请输入搜索关键词'
}
```

### 日期选择字段

```javascript
{
  name: 'start_date',
  label: '开始日期',
  type: 'date',
  required: true
}
```

### 数字输入字段

```javascript
{
  name: 'amount',
  label: '金额',
  type: 'number',
  required: true,
  placeholder: '请输入金额'
}
```

### 文本域字段

```javascript
{
  name: 'description',
  label: '描述',
  type: 'textarea',
  required: false,
  placeholder: '请输入详细描述'
}
```

---

## 💡 实际案例

### 案例 1：GitHub Trending（简单执行）

**n8n 配置：**
- Webhook 触发器（GET）
- 获取 GitHub API 数据
- 生成 HTML 报告

**前端配置：**
```javascript
'My workflow 2': {
  type: WORKFLOW_TYPES.SIMPLE,
  description: '获取 GitHub 本周热门项目',
  executeLabel: '获取热门项目',
  resultType: 'html'
}
```

**用户体验：**
1. 点击"执行"按钮
2. 确认执行
3. 查看 HTML 结果

---

### 案例 2：Azure 账单对比（文件上传）

**n8n 配置：**
- Webhook 触发器（POST）
- 接收两个 Excel 文件
- 对比分析
- 返回 Excel 报告

**前端配置：**
```javascript
'Azure账单对比': {
  type: WORKFLOW_TYPES.UPLOAD_FILES,
  description: '上传本月和上月的 Azure 账单进行对比',
  executeLabel: '开始对比',
  resultType: 'download',
  fields: [
    {
      name: 'current_month',
      label: '本月账单',
      type: 'file',
      accept: '.xlsx,.xls',
      required: true
    },
    {
      name: 'last_month',
      label: '上月账单',
      type: 'file',
      accept: '.xlsx,.xls',
      required: true
    }
  ]
}
```

**用户体验：**
1. 点击"执行"按钮
2. 上传两个 Excel 文件
3. 点击"开始对比"
4. 自动下载对比报告

---

### 案例 3：数据查询（表单输入）

**n8n 配置：**
- Webhook 触发器（POST）
- 接收查询条件
- 查询数据库
- 返回 JSON 结果

**前端配置：**
```javascript
'数据查询': {
  type: WORKFLOW_TYPES.FORM_INPUT,
  description: '根据日期和关键词查询数据',
  executeLabel: '查询',
  resultType: 'json',
  fields: [
    {
      name: 'startDate',
      label: '开始日期',
      type: 'date',
      required: true
    },
    {
      name: 'endDate',
      label: '结束日期',
      type: 'date',
      required: true
    },
    {
      name: 'keyword',
      label: '关键词',
      type: 'text',
      required: false,
      placeholder: '可选'
    }
  ]
}
```

**用户体验：**
1. 点击"执行"按钮
2. 填写日期和关键词
3. 点击"查询"
4. 查看 JSON 结果

---

## 🔧 高级功能

### 自定义结果处理

如果需要特殊的结果处理逻辑，可以在 `resultType` 中指定：

- `html` - 在 iframe 中渲染 HTML
- `json` - 格式化显示 JSON
- `download` - 自动触发文件下载
- `auto` - 自动检测（默认）

### 添加新的工作流类型

如需添加新类型，编辑 `workflow-configs.js`：

```javascript
export const WORKFLOW_TYPES = {
  SIMPLE: 'simple',
  UPLOAD_FILES: 'upload_files',
  FORM_INPUT: 'form_input',
  YOUR_NEW_TYPE: 'your_new_type'  // 添加新类型
}
```

然后在 `DynamicExecuteDialog.vue` 中添加对应的 UI 逻辑。

---

## ✅ 最佳实践

### 1. 工作流命名
- 使用清晰、描述性的名称
- 中英文都可以
- 避免特殊字符

### 2. Webhook 配置
- 始终使用 Production URL（激活工作流后可用）
- 设置合适的超时时间
- 添加 CORS 头（如果需要）

### 3. 字段设计
- 必填字段尽量少
- 提供清晰的 placeholder
- 添加字段描述

### 4. 错误处理
- 在 n8n 中添加错误处理节点
- 返回友好的错误信息
- 记录执行日志

---

## 🎨 UI 自定义

所有工作流卡片会自动显示：
- 工作流类型标签（颜色编码）
- 激活状态
- 更新时间
- 节点数量

无需手动配置！

---

## 🚨 常见问题

### Q: 新工作流没有显示？
**A:** 检查：
1. n8n 中是否已保存工作流
2. 前端是否已刷新
3. 点击"刷新"按钮获取最新列表

### Q: 执行按钮点击没反应？
**A:** 检查：
1. 工作流是否已激活（n8n 中）
2. Webhook 是否配置正确
3. 浏览器控制台是否有错误

### Q: 如何调试工作流？
**A:**
1. 打开浏览器控制台（F12）
2. 点击执行工作流
3. 查看网络请求和控制台日志
4. 在 n8n 中查看执行历史

### Q: 文件上传失败？
**A:** 检查：
1. 文件大小是否超限
2. 文件类型是否正确
3. n8n Webhook 是否接受 POST 请求
4. Content-Type 是否为 multipart/form-data

---

## 📚 参考资料

- [n8n 官方文档](https://docs.n8n.io/)
- [Element Plus 组件库](https://element-plus.org/)
- [Vue 3 文档](https://vuejs.org/)

---

## 🎉 总结

使用这个架构，你可以：

✅ **快速集成** - 新工作流只需配置，无需编码  
✅ **类型安全** - 明确的字段定义和验证  
✅ **用户友好** - 自动生成适配的 UI  
✅ **易于维护** - 配置集中管理  
✅ **可扩展** - 轻松添加新类型和功能  

**开始创建你的工作流吧！** 🚀
