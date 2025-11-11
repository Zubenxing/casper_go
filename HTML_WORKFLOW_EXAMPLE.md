# 📝 HTML 工作流返回示例

## ✨ 功能说明

前端现在支持 **自动识别并渲染 HTML 内容**！

当你的 n8n 工作流返回 HTML 格式数据时：
- ✅ **自动检测** HTML 内容
- 🎨 **直接渲染** 在页面中
- 📋 **一键复制** HTML 代码
- 🪟 **新窗口查看** 完整页面

---

## 🔧 n8n 工作流配置

### 方式 1: 在 Code 节点返回 HTML

在你的工作流最后添加一个 **Code** 节点：

```javascript
// 生成 HTML 内容
const html = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>GitHub 热门项目</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      padding: 20px;
      background: #f5f7fa;
    }
    .container {
      max-width: 1200px;
      margin: 0 auto;
      background: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 12px rgba(0,0,0,0.1);
    }
    h1 {
      color: #409eff;
      border-bottom: 2px solid #409eff;
      padding-bottom: 10px;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 20px;
    }
    th, td {
      padding: 12px;
      text-align: left;
      border-bottom: 1px solid #ddd;
    }
    th {
      background-color: #409eff;
      color: white;
    }
    tr:hover {
      background-color: #f5f7fa;
    }
    .stars {
      color: #f39c12;
      font-weight: bold;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>🔥 GitHub 每周热门项目 Top 20</h1>
    <table>
      <thead>
        <tr>
          <th>#</th>
          <th>项目名称</th>
          <th>Stars</th>
          <th>语言</th>
          <th>描述</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>1</td>
          <td><a href="https://github.com/awesome/project">awesome-project</a></td>
          <td class="stars">⭐ 12,345</td>
          <td>Go</td>
          <td>An awesome project for developers</td>
        </tr>
        <!-- 更多数据... -->
      </tbody>
    </table>
    <p style="margin-top: 20px; color: #909399; text-align: center;">
      数据更新时间: ${new Date().toLocaleString('zh-CN')}
    </p>
  </div>
</body>
</html>
`;

// 返回 HTML
return [
  {
    json: {
      html: html,  // ← 关键：使用 'html' 字段
      timestamp: new Date().toISOString()
    }
  }
];
```

**关键点：**
- 返回对象的 `json.html` 字段
- 或者使用 `json.output` 或 `json.result`

---

### 方式 2: HTTP Request 返回 HTML

如果你的工作流调用外部 API 返回 HTML：

1. **HTTP Request** 节点配置：
   - URL: 你的 API 地址
   - Response Format: `String`（不要用 JSON）

2. **Code** 节点处理：
```javascript
// 从上一步获取 HTML
const html = $input.all()[0].json.data

// 返回
return [
  {
    json: {
      html: html
    }
  }
];
```

---

### 方式 3: 动态生成报表

```javascript
// 获取数据（从上一个节点）
const projects = $input.all()[0].json.projects || []

// 生成 HTML 表格
let tableRows = projects.map((project, index) => `
  <tr>
    <td>${index + 1}</td>
    <td><a href="${project.url}" target="_blank">${project.name}</a></td>
    <td class="stars">⭐ ${project.stars.toLocaleString()}</td>
    <td>${project.language}</td>
    <td>${project.description || '-'}</td>
  </tr>
`).join('')

// 完整 HTML
const html = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>GitHub 热门项目报告</title>
  <style>
    body { 
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      margin: 0;
      padding: 20px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }
    .container {
      max-width: 1400px;
      margin: 0 auto;
      background: white;
      padding: 30px;
      border-radius: 12px;
      box-shadow: 0 10px 40px rgba(0,0,0,0.2);
    }
    h1 {
      color: #667eea;
      font-size: 32px;
      margin-bottom: 10px;
    }
    .subtitle {
      color: #999;
      margin-bottom: 30px;
    }
    table {
      width: 100%;
      border-collapse: collapse;
    }
    th {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      padding: 15px;
      text-align: left;
      font-weight: 600;
    }
    td {
      padding: 15px;
      border-bottom: 1px solid #eee;
    }
    tr:hover {
      background-color: #f8f9fa;
    }
    .stars {
      color: #f39c12;
      font-weight: bold;
      font-size: 16px;
    }
    a {
      color: #667eea;
      text-decoration: none;
      font-weight: 600;
    }
    a:hover {
      text-decoration: underline;
    }
    .footer {
      margin-top: 30px;
      text-align: center;
      color: #999;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>🚀 GitHub 热门项目周报</h1>
    <p class="subtitle">本周最受关注的 ${projects.length} 个开源项目</p>
    
    <table>
      <thead>
        <tr>
          <th style="width: 50px">#</th>
          <th>项目</th>
          <th style="width: 120px">Stars</th>
          <th style="width: 100px">语言</th>
          <th>描述</th>
        </tr>
      </thead>
      <tbody>
        ${tableRows}
      </tbody>
    </table>
    
    <div class="footer">
      <p>📅 生成时间: ${new Date().toLocaleString('zh-CN')}</p>
      <p>💡 数据来源: GitHub Trending API</p>
    </div>
  </div>
</body>
</html>
`;

return [{ json: { html: html } }];
```

---

## 🎯 前端展示效果

执行工作流后，前端会：

1. **自动检测** HTML 内容
2. **显示标签**：`HTML` 或 `JSON`
3. **直接渲染** HTML（支持样式、表格、图片等）
4. **提供操作**：
   - 📋 **复制 HTML** - 复制完整 HTML 代码
   - 🪟 **新窗口查看** - 在新标签页中打开
   - 🔍 **在 n8n 中查看详情** - 查看执行日志

---

## 📋 支持的 HTML 元素

前端会自动美化以下元素：

| 元素 | 样式 |
|------|------|
| `<table>` | 边框、间距、悬停效果 |
| `<img>` | 响应式宽度 |
| `<pre>` | 代码块背景色 |
| `<code>` | 内联代码样式 |
| 所有元素 | 自动适应容器宽度 |

---

## 💡 使用场景

### 1. 数据报表
生成可视化的数据统计报表，直接展示给用户

### 2. 邮件预览
生成邮件 HTML，在发送前预览效果

### 3. 爬虫结果
网页抓取后，保留原始样式展示

### 4. 动态页面
根据数据生成动态 HTML 页面

### 5. 文档生成
自动生成格式化的技术文档

---

## 🔍 调试技巧

### 1. 检查返回格式

在 n8n 中点击节点查看输出：
```json
{
  "json": {
    "html": "<html>...</html>",  // ✅ 正确
    "timestamp": "2025-11-10T08:00:00Z"
  }
}
```

### 2. 使用 console.log

在 Code 节点中调试：
```javascript
const html = generateHTML()
console.log('HTML length:', html.length)
console.log('Has table:', html.includes('<table>'))
return [{ json: { html } }]
```

### 3. 前端检查

打开浏览器开发者工具，查看：
- Network 标签：查看 API 返回的原始数据
- Console 标签：查看是否有错误信息

---

## ⚠️ 注意事项

### 1. XSS 安全
- 前端使用 `v-html` 渲染，请确保 HTML 内容可信
- 不要渲染用户输入的未经处理的 HTML

### 2. 样式隔离
- 使用内联样式或 `<style>` 标签
- 避免使用全局 CSS 类名冲突

### 3. 图片路径
- 使用完整 URL：`https://example.com/image.png`
- 或使用 Base64：`data:image/png;base64,...`

### 4. 脚本限制
- `<script>` 标签会被浏览器执行
- 谨慎使用 JavaScript

---

## 🎉 示例效果

当你执行工作流后，页面会显示：

```
┌─────────────────────────────────────┐
│ ✅ 执行成功                  [关闭] │
│ 工作流：My workflow                 │
│ 执行ID：12345                       │
│ 开始时间：2025-11-10 16:00         │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ 📦 返回数据          [HTML]         │
│ ┌───────────────────────────────┐   │
│ │                               │   │
│ │  🚀 GitHub 热门项目周报       │   │
│ │                               │   │
│ │  本周最受关注的 20 个项目     │   │
│ │                               │   │
│ │  # | 项目 | Stars | 语言      │   │
│ │  ────────────────────────────  │   │
│ │  1 | awesome-go | ⭐12K | Go │   │
│ │  2 | vue | ⭐200K | JS       │   │
│ │  ...                          │   │
│ │                               │   │
│ └───────────────────────────────┘   │
│                                     │
│  [📋 复制 HTML] [🪟 新窗口查看]     │
│  [🔍 在n8n中查看详情]               │
└─────────────────────────────────────┘
```

**HTML 内容会直接渲染，表格、样式、颜色等都会正常显示！** 🎨

---

## 🚀 快速测试

1. 在 n8n 中创建简单的工作流
2. 添加 Code 节点，返回示例 HTML
3. 保存并在前端点击"执行"
4. 查看渲染效果！

就是这么简单！享受 HTML 渲染的便利吧！✨
