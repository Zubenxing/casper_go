# 路由模块化配置说明

## 📁 目录结构

```
routes/
├── index.js          # 路由聚合文件
├── dashboard.js      # 首页路由
├── certificates.js   # 证书监控路由
├── profile.js        # 个人中心路由
└── README.md         # 说明文档
```

## 📝 新增路由步骤

### 1. 创建路由配置文件

在 `routes/` 目录下创建新的路由文件，例如 `example.js`：

```javascript
/**
 * 示例模块路由配置
 */
export default {
  path: '/example',
  name: 'Example',
  component: () => import('@/features/example/Example.vue'),
  meta: {
    title: '示例页面',      // 页面标题
    icon: 'Setting',       // 菜单图标（Element Plus 图标名）
    showInMenu: true,      // 是否在侧边栏菜单显示
    order: 4,              // 菜单排序（数字越小越靠前）
    
    // 可选：标签栏配置
    hasTabBar: false,      // 是否显示顶部标签栏
    tabs: []               // 标签栏配置（如果 hasTabBar 为 true）
  }
}
```

### 2. 在 `routes/index.js` 中导入

```javascript
import example from './example'

export const moduleRoutes = [
  dashboard,
  certificates,
  profile,
  example  // 添加新路由
].sort((a, b) => (a.meta?.order || 999) - (b.meta?.order || 999))
```

### 3. 创建对应的 Vue 组件

在 `frontend/src/features/` 下创建功能模块：

```
features/
└── example/
    ├── Example.vue
    ├── api.js          # 可选：API 定义
    └── store.js        # 可选：状态管理
```

## 🎨 特殊功能配置

### 标签栏（Tab Bar）

如需在页面顶部显示标签栏（如证书监控页面），配置：

```javascript
meta: {
  hasTabBar: true,
  tabs: [
    { name: 'list', label: '列表' },
    { name: 'create', label: '创建' }
  ]
}
```

### 菜单图标

使用 Element Plus 图标库的图标名：
- `HomeFilled` - 首页
- `Document` - 文档
- `User` - 用户
- `Setting` - 设置
- 更多图标查看：https://element-plus.org/zh-CN/component/icon.html

## 🔧 配置项说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `path` | string | ✅ | 路由路径 |
| `name` | string | ✅ | 路由名称（唯一） |
| `component` | function | ✅ | 组件导入函数 |
| `meta.title` | string | ✅ | 页面标题 |
| `meta.icon` | string | ✅ | 菜单图标名 |
| `meta.showInMenu` | boolean | ✅ | 是否显示在菜单 |
| `meta.order` | number | ❌ | 菜单排序（默认 999） |
| `meta.hasTabBar` | boolean | ❌ | 是否有标签栏 |
| `meta.tabs` | array | ❌ | 标签栏配置 |

## 🌟 优势

1. **模块化**：每个功能一个文件，易于维护
2. **可扩展**：新增路由无需修改核心代码
3. **配置驱动**：路由配置自动生成菜单和导航
4. **类型安全**：统一的配置结构
5. **灵活性**：支持特殊功能（如标签栏）

## 📖 示例

参考 `certificates.js` 查看完整的标签栏配置示例。

