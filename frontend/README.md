# Casper Platform 前端

基于 Vue 3 + Vite + Element Plus 构建的现代化证书监控平台前端。

## 技术栈

- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **UI 组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP**: Axios

## 项目结构（按功能模块组织）

```
frontend/src/
├── features/          # 功能模块（每个功能独立）
│   ├── auth/         # 认证模块
│   │   ├── Login.vue      # 登录页面
│   │   ├── Profile.vue    # 个人中心页面
│   │   ├── api.js         # 认证相关API
│   │   └── store.js       # 认证状态管理
│   ├── dashboard/    # 首页模块
│   │   └── Dashboard.vue  # 首页
│   └── certificates/ # 证书监控模块
│       ├── Certificates.vue  # 证书列表页面
│       ├── api.js            # 证书相关API
│       └── components/       # 证书相关组件（可扩展）
├── core/             # 核心公共功能
│   ├── layouts/      # 布局组件
│   │   └── MainLayout.vue
│   ├── router/       # 路由配置
│   │   └── index.js
│   ├── api/          # API 配置
│   │   ├── request.js     # Axios 实例
│   │   └── index.js       # API 统一导出
│   └── utils/        # 工具函数
│       └── format.js      # 格式化工具
├── App.vue           # 根组件
└── main.js           # 入口文件
```

## 目录设计原则

### features/ - 功能模块
每个业务功能独立一个目录，包含：
- **页面组件** (.vue 文件)
- **API 接口** (api.js)
- **状态管理** (store.js，可选)
- **子组件** (components/，可选)

**优点：**
- ✅ 功能独立，代码内聚
- ✅ 易于维护和扩展
- ✅ 新增功能只需创建新目录
- ✅ 删除功能直接删除目录

### core/ - 核心公共
所有功能模块共享的代码：
- **layouts/** - 布局组件
- **router/** - 路由配置
- **api/** - HTTP 请求配置
- **utils/** - 工具函数

## 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

访问：http://localhost:3000

### 3. 构建生产版本

```bash
npm run build
```

## 添加新功能

假设要添加一个"通知管理"功能：

### 1. 创建功能目录

```bash
mkdir src/features/notifications
```

### 2. 创建文件

```
features/notifications/
├── Notifications.vue  # 页面
├── api.js            # API 接口
└── store.js          # 状态管理（可选）
```

### 3. 编写 API

```javascript
// features/notifications/api.js
import request from '@/core/api/request'

export default {
  getList: () => request.get('/notifications'),
  markRead: (id) => request.put(`/notifications/${id}/read`)
}
```

### 4. 注册 API

```javascript
// core/api/index.js
import notificationsApi from '@/features/notifications/api'

export default {
  // ...
  notifications: notificationsApi
}
```

### 5. 添加路由

```javascript
// core/router/index.js
{
  path: '/notifications',
  name: 'Notifications',
  component: () => import('@/features/notifications/Notifications.vue'),
  meta: { title: '通知管理' }
}
```

### 6. 添加菜单

在 `core/layouts/MainLayout.vue` 的菜单中添加入口。

就这么简单！

## 开发指南

### API 调用

```javascript
import api from '@/core/api'

// 调用
const response = await api.certificates.getList(1, 10)
```

### 使用工具函数

```javascript
import { formatTime, getDaysLeftType } from '@/core/utils/format'

const time = formatTime(new Date())
const type = getDaysLeftType(30) // 返回 'warning'
```

### 使用状态管理

```javascript
import { useAuthStore } from '@/features/auth/store'

const authStore = useAuthStore()
console.log(authStore.user)
```

## 配置说明

### API 代理

开发环境已配置代理（`vite.config.js`）：
```javascript
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}
```

### 生产环境

修改 `core/api/request.js` 的 `baseURL`：
```javascript
baseURL: 'https://your-domain.com/api'
```

## 响应式设计

所有页面都使用 Element Plus 的栅格系统：
- **xs**: <768px (手机)
- **sm**: ≥768px (平板)
- **md**: ≥992px (桌面)
- **lg**: ≥1200px (大屏)

## 默认账号

- 用户名: `admin`
- 密码: `admin123`

## 更多

查看 `FRONTEND_API.md` 了解完整的 API 接口文档。
