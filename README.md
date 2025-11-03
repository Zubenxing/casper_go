# Casper - SSL证书监控平台

一个基于 Go + Vue 3 的现代化 SSL 证书监控和管理平台。

## ✨ 主要功能

### 🔐 证书监控
- SSL 证书到期监控
- 自动定时检查证书状态
- 证书信息详细展示
- 到期提醒（正常/警告/过期）
- 批量导入/导出证书

### 🛠️ 证书工具
- **CSR 生成器**：在线生成证书签名请求和私钥
  - 支持 RSA 和 ECDSA 算法
  - 可自定义密钥长度
  - 支持多域名（SAN）
- **CSR 验证**：验证 CSR 文件有效性
- **证书验证**：验证证书和私钥配对关系

### 🔑 密码管理
- 安全存储网站账户密码
- AES-256-GCM 加密
- 分类管理
- 收藏功能
- 快速搜索和筛选

## 🚀 技术栈

### 后端
- **Go 1.21+**
- **Gin** - Web 框架
- **GORM** - ORM
- **MySQL** - 数据库
- **JWT** - 身份认证
- **Swagger** - API 文档

### 前端
- **Vue 3** - 前端框架
- **Element Plus** - UI 组件库
- **Vite** - 构建工具
- **Vue Router** - 路由管理
- **Pinia** - 状态管理

## 📦 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 1. 克隆项目

```bash
git clone https://github.com/your-username/casper_go.git
cd casper_go
```

### 2. 配置数据库

```bash
# 复制配置文件
cp config/database.yaml.example config/database.yaml
cp config/site_info.yaml.example config/site_info.yaml

# 编辑配置文件，填入你的数据库信息
# config/database.yaml
```

### 3. 执行数据库迁移

```bash
# Windows
.\migrate.ps1 up

# Linux/Mac
./migrate.sh up

# 或使用 Makefile
make migrate-up
```

### 4. 启动后端

```bash
# 开发模式
go run cmd/server/main.go

# 或编译后运行
go build -o casper_server cmd/server/main.go
./casper_server
```

后端将在 `http://localhost:8088` 启动

### 5. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 `http://localhost:3001` 启动

### 6. 登录系统

默认管理员账户：
- 用户名：`admin`
- 密码：`admin123`

⚠️ **首次登录后请立即修改密码！**

## 📖 文档

- [数据库迁移指南](docs/MIGRATION_GUIDE.md)
- [性能优化总结](docs/OPTIMIZATION_SUMMARY.md)
- [证书工具 API](docs/CERTIFICATE_TOOLS_API.md)
- [密码管理 API](docs/PASSWORD_API.md)
- [前端 API 文档](FRONTEND_API.md)

## 🔧 开发命令

### 后端

```bash
# 安装依赖
make install

# 生成 Swagger 文档
make swagger

# 运行项目
make run

# 编译项目
make build

# 运行测试
make test

# 代码检查
make lint
```

### 前端

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 预览构建结果
npm run preview
```

### 数据库迁移

```bash
# 查看迁移状态
make migrate-status

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down
```

## 📁 项目结构

```
casper_go/
├── cmd/                    # 命令行入口
│   ├── server/            # 服务器主程序
│   └── migrate/           # 数据库迁移工具
├── core/                   # 核心功能
│   ├── database/          # 数据库连接和迁移
│   ├── logger/            # 日志系统
│   ├── middleware/        # 中间件
│   ├── response/          # 响应封装
│   └── utils/             # 工具函数
├── features/               # 业务功能模块
│   ├── auth/              # 用户认证
│   ├── certificate/       # 证书监控
│   └── password/          # 密码管理
├── migrations/             # 数据库迁移文件
├── router/                 # 路由配置
├── config/                 # 配置文件
├── frontend/               # 前端项目
│   ├── src/
│   │   ├── core/          # 核心功能（路由、API等）
│   │   └── features/      # 业务功能组件
│   └── package.json
├── docs/                   # 文档
└── README.md
```

## 🎯 核心特性

### 并发证书检查
- 10个 goroutine 并发处理
- 性能提升 8+ 倍
- 实时进度反馈
- 详细统计信息

### 数据库迁移系统
- 类似 Laravel 的 migration 机制
- 支持 batch 批次管理
- 事务保护，失败自动回滚
- 可追踪、可回滚

### 安全加密
- AES-256-GCM 加密算法
- 每个密码独立加密密钥
- JWT 身份认证
- 密码加密存储

### 性能优化
- 数据库索引优化（8个关键索引）
- 查询性能提升 40+ 倍
- 前端虚拟滚动
- 请求防抖和节流

## 🔒 安全建议

1. **修改默认密码**：首次登录后立即修改 admin 账户密码
2. **配置文件安全**：不要将 `config/database.yaml` 提交到版本控制
3. **JWT 密钥**：生产环境使用强随机密钥
4. **HTTPS**：生产环境启用 HTTPS
5. **定期更新**：保持依赖包更新

## 📝 更新日志

### v1.0.0 (2025-11-03)
- ✅ SSL 证书监控功能
- ✅ 证书工具（CSR 生成/验证、证书验证）
- ✅ 密码管理功能
- ✅ 数据库迁移系统
- ✅ 并发优化（性能提升 8+ 倍）
- ✅ 数据库索引优化（性能提升 40+ 倍）
- ✅ 前端用户体验优化

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

[MIT License](LICENSE)

## 👨‍💻 作者

Your Name

## 🙏 鸣谢

- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [Vue 3](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
