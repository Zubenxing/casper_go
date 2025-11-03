# 证书工具前端开发文档

## 📋 概述

证书工具前端提供三个功能页面，与后端 API 完全对接：
1. **生成 CSR** - 在线生成证书签名请求和私钥
2. **验证 CSR** - 验证 CSR 文件的合法性
3. **验证证书** - 验证 SSL 证书文件

## 🗂️ 文件结构

```
frontend/src/
├── core/
│   ├── router/
│   │   └── routes/
│   │       └── cert-tools.js        # 路由配置 ✨新增
│   └── layouts/
│       └── MainLayout.vue            # 更新：支持动态 tab
└── features/
    └── cert-tools/                   # 证书工具模块 ✨新增
        ├── api.js                    # API 定义
        ├── CertTools.vue             # 容器组件
        ├── GenerateCSR.vue           # CSR 生成页面
        ├── ValidateCSR.vue           # CSR 验证页面
        └── ValidateCert.vue          # 证书验证页面
```

---

## 🎨 功能详情

### 1. 生成 CSR (`GenerateCSR.vue`)

#### 功能特性
- ✅ 完整的表单验证
- ✅ 支持 RSA (2048/4096) 和 ECDSA (P-256/P-384)
- ✅ 动态密钥长度选择
- ✅ SAN (Subject Alternative Names) 支持
- ✅ 一键下载私钥和 CSR 文件
- ✅ 安全提示和下一步操作指引

#### 表单字段
- **必填**：通用名 (CN)、密钥算法、密钥长度
- **可选**：国家、省份、城市、组织、部门、邮箱、DNS 名称

#### 用户体验
- 表单提示说明
- 实时验证
- 密钥算法切换自动调整密钥长度选项
- 生成成功后自动滚动到结果区域
- 私钥和 CSR 独立下载按钮

---

### 2. 验证 CSR (`ValidateCSR.vue`)

#### 功能特性
- ✅ 支持粘贴 PEM 格式内容
- ✅ 支持拖拽上传文件
- ✅ 支持 .csr, .pem, .txt 文件
- ✅ 详细的验证结果展示
- ✅ 错误提示

#### 展示信息
- 验证状态（成功/失败）
- 主题信息（CN, C, ST, L, O, OU）
- SAN 信息（DNS 名称、邮箱地址）
- 公钥信息（算法、长度）
- 签名算法

---

### 3. 验证证书 (`ValidateCert.vue`)

#### 功能特性
- ✅ 支持粘贴 PEM 格式证书
- ✅ 支持拖拽上传文件
- ✅ 支持 .crt, .cer, .pem, .txt 文件
- ✅ 全面的证书信息展示
- ✅ 证书状态智能判断

#### 展示信息
**证书状态**
- 有效 / 即将过期 / 接近过期 / 已过期 / 尚未生效
- 剩余天数

**基本信息**
- 版本、序列号
- 通用名、组织、国家
- 是否为 CA 证书

**颁发者和主题**
- 完整的 DN (Distinguished Name)

**有效期**
- 生效时间
- 过期时间

**密钥信息**
- 公钥算法（RSA/ECDSA）
- 公钥长度
- 签名算法

**SAN 信息**
- DNS 名称列表
- IP 地址列表
- 邮箱地址列表

**密钥用途**
- Key Usage
- Extended Key Usage

---

## 🎯 路由配置

```javascript
{
  path: '/cert-tools',
  name: 'CertTools',
  component: () => import('@/features/cert-tools/CertTools.vue'),
  meta: {
    title: '证书工具',
    icon: 'Tools',
    showInMenu: true,
    order: 3,
    hasTabBar: true,
    tabs: [
      { name: 'generate-csr', label: '生成 CSR' },
      { name: 'validate-csr', label: '验证 CSR' },
      { name: 'validate-cert', label: '验证证书' }
    ]
  }
}
```

---

## 🔌 API 集成

### API 定义 (`api.js`)

```javascript
import http from '@/core/api/request'

export default {
  // 生成 CSR
  generateCSR(data) {
    return http.post('/certificates/tools/generate-csr', data)
  },
  
  // 验证 CSR
  validateCSR(data) {
    return http.post('/certificates/tools/validate-csr', data)
  },
  
  // 验证证书
  validateCert(data) {
    return http.post('/certificates/tools/validate-cert', data)
  }
}
```

### 请求/响应示例

参见 [CERTIFICATE_TOOLS_API.md](./CERTIFICATE_TOOLS_API.md)

---

## 🎨 UI 设计

### 设计原则
1. **简洁直观** - 清晰的表单布局和结果展示
2. **提示友好** - 每个字段都有说明
3. **反馈及时** - Loading 状态、成功/失败提示
4. **操作便捷** - 一键下载、拖拽上传

### 色彩规范
- 成功：`#67c23a`
- 警告：`#e6a23c`
- 危险：`#f56c6c`
- 信息：`#409eff`

### 组件使用
- Element Plus 组件库
- 图标：`@element-plus/icons-vue`
- 卡片布局：`el-card`
- 表单：`el-form` + `el-form-item`
- 标签：`el-tag`
- 上传：`el-upload` (拖拽模式)

---

## 🔄 状态管理

### Tab 状态
- 通过 MainLayout 的 `provide/inject` 机制管理
- 自动根据路由配置初始化
- 路由切换时自动重置

### 表单状态
- 使用 Vue 3 `reactive` 管理表单数据
- `ref` 管理 Loading 和结果状态

---

## ✨ 用户体验优化

### GenerateCSR
1. **智能默认值**：选择 RSA 自动设置 2048 位，选择 ECDSA 自动设置 P-256
2. **安全提示**：明确提示私钥仅显示一次
3. **操作指引**：生成成功后显示详细的下一步操作

### ValidateCSR / ValidateCert
1. **多种输入方式**：粘贴文本 + 拖拽上传 + 点击选择
2. **格式提示**：明确说明支持的文件格式
3. **结果分组**：信息按类别组织，易于阅读
4. **智能判断**：证书验证自动判断状态（过期、即将过期等）

---

## 📱 响应式设计

所有页面支持响应式布局：
- 表单：单列布局，适配小屏幕
- 结果展示：使用 Grid 布局，自动适应
- 最大宽度：1200px，居中显示

---

## 🚀 使用流程

### 生成 CSR
1. 填写表单（至少填写通用名）
2. 选择密钥算法和长度
3. 点击"生成 CSR"
4. 下载私钥文件（重要！）
5. 下载 CSR 文件
6. 将 CSR 提交给 CA

### 验证 CSR
1. 粘贴 CSR 内容或拖拽文件上传
2. 点击"验证 CSR"
3. 查看验证结果

### 验证证书
1. 粘贴证书内容或拖拽文件上传
2. 点击"验证证书"
3. 查看证书详细信息和状态

---

## 🧪 测试建议

1. **生成 CSR**
   - 测试 RSA 2048/4096
   - 测试 ECDSA P-256/P-384
   - 测试带 SAN 的 CSR
   - 验证下载功能

2. **验证 CSR**
   - 测试有效的 CSR
   - 测试无效的 PEM 格式
   - 测试文件上传

3. **验证证书**
   - 测试有效证书
   - 测试过期证书
   - 测试包含 SAN 的证书
   - 测试 CA 证书

---

## 📊 性能优化

1. **懒加载**：所有页面组件使用动态 import
2. **防抖**：文件上传和表单提交有 Loading 状态防止重复
3. **内存管理**：文件读取后正确释放资源

---

## 🔧 维护指南

### 新增字段
如需在生成 CSR 时新增字段：
1. 在 `GenerateCSR.vue` 的 `form` 中添加字段
2. 在表单中添加对应的 `el-form-item`
3. 在 `handleGenerate` 中包含该字段
4. 更新表单验证规则（如需要）

### 调整样式
1. 所有样式都使用 `scoped`，不会影响其他组件
2. 使用 Element Plus 的主题变量保持一致性
3. 自定义样式使用类名命名规范

---

## 🎉 完成状态

- ✅ 路由配置
- ✅ API 集成
- ✅ 生成 CSR 页面
- ✅ 验证 CSR 页面
- ✅ 验证证书页面
- ✅ Tab 导航
- ✅ 响应式布局
- ✅ 错误处理
- ✅ 用户体验优化

**前端证书工具开发完成！** 🎊

