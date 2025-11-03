# API 汇总文档

## 🎯 后端 API 总览

### 认证相关
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/refresh` - 刷新 Token
- `GET /api/profile` - 获取用户信息（需认证）
- `POST /api/logout` - 退出登录（需认证）

### 证书监控
- `POST /api/certificates` - 添加证书监控（需认证）
- `GET /api/certificates` - 获取证书列表（需认证）
- `GET /api/certificates/:id` - 获取证书详情（需认证）
- `PUT /api/certificates/:id` - 更新证书信息（需认证）
- `PATCH /api/certificates/:id/info` - 更新客户名和备注（需认证）
- `DELETE /api/certificates/:id` - 删除证书监控（需认证）
- `POST /api/certificates/check-all` - 检查所有证书（需认证）

### 证书工具（公开访问）
- `POST /api/certificates/tools/generate-csr` - 生成 CSR
- `POST /api/certificates/tools/validate-csr` - 验证 CSR
- `POST /api/certificates/tools/validate-cert` - 验证证书

---

## 📊 功能完成状态

### ✅ 已完成

#### 后端
- [x] 用户认证系统
- [x] 证书监控 CRUD
- [x] 证书自动检查
- [x] CSR 在线生成（RSA & ECDSA）
- [x] CSR 在线验证
- [x] SSL 证书在线验证
- [x] 客户名和备注功能
- [x] 证书导入导出

#### 前端
- [x] 路由模块化重构
- [x] 证书列表页面
- [x] 证书更新页面（UI）
- [x] 多格式导出（JSON/CSV/Excel/TXT）
- [x] 多格式导入（JSON/CSV/TXT）
- [x] 拖拽上传支持
- [x] 自动刷新（每天 00:05）
- [x] 编辑证书信息

### 🚧 待开发

#### 前端
- [ ] CSR 生成页面
- [ ] CSR 验证页面
- [ ] 证书验证页面
- [ ] 批量更新证书功能
- [ ] 即将过期证书筛选
- [ ] 已过期证书筛选

---

## 🗂️ 项目结构

```
casper_go/
├── backend/
│   ├── features/
│   │   ├── auth/           # 认证功能
│   │   └── certificate/    # 证书相关
│   │       ├── api.go      # API 处理器（新增 3 个工具 API）
│   │       ├── checker.go  # 证书检查器
│   │       ├── service.go  # 业务逻辑
│   │       ├── models.go   # 数据模型
│   │       ├── csr.go      # CSR 生成和验证 ✨新增
│   │       └── cert_validator.go  # 证书验证 ✨新增
│   ├── core/
│   │   ├── api/
│   │   ├── database/
│   │   ├── logger/
│   │   └── middleware/
│   ├── router/
│   │   └── router.go       # 路由配置（新增工具路由）
│   └── config/
├── frontend/
│   └── src/
│       ├── features/
│       │   ├── auth/
│       │   ├── dashboard/
│       │   └── certificates/
│       │       ├── Certificates.vue     # 容器组件
│       │       ├── CertificateList.vue  # 列表页
│       │       └── CertUpdate.vue       # 更新页
│       ├── core/
│       │   ├── router/
│       │   │   └── routes/  # 路由模块化 ✨新增
│       │   ├── layouts/
│       │   ├── api/
│       │   └── config/
│       │       └── menu.js  # 菜单配置 ✨新增
│       └── utils/
├── docs/
│   ├── CERTIFICATE_TOOLS_API.md  # 证书工具 API 文档 ✨新增
│   └── API_SUMMARY.md             # API 汇总文档 ✨新增
└── test_data/  # 测试数据 ✨新增
```

---

## 🎨 前端开发指南

### 下一步：开发证书工具页面

#### 1. 创建路由
在 `frontend/src/core/router/routes/` 下创建 `cert-tools.js`：

```javascript
export default {
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

#### 2. 创建 API 定义
在 `frontend/src/features/cert-tools/api.js`：

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

#### 3. 创建组件
- `CertTools.vue` - 容器组件（类似 `Certificates.vue`）
- `GenerateCSR.vue` - CSR 生成页面
- `ValidateCSR.vue` - CSR 验证页面
- `ValidateCert.vue` - 证书验证页面

---

## 🔐 测试账号

- **用户名**: `admin`
- **密码**: `admin123`

---

## 📚 相关文档

- [证书工具 API 文档](./CERTIFICATE_TOOLS_API.md)
- [路由模块化说明](../frontend/src/core/router/routes/README.md)
- [Swagger 文档](http://localhost:8080/swagger/index.html)

