# 密码管理 API 文档

## 📚 概述

密码管理功能用于安全存储和管理用户的账户密码信息。所有密码使用 AES-256-GCM 加密存储。

## 🔐 安全说明

- ✅ 密码使用 AES-256-GCM 加密存储
- ✅ 密码不会在列表和详情接口中返回
- ✅ 查看密码需要单独调用接口
- ✅ 所有接口需要JWT认证
- ✅ 用户只能访问自己的数据

## 📡 API 端点

### 1. 创建账户

**POST** `/api/passwords`

#### 请求体
```json
{
  "title": "Google账户",
  "url": "https://accounts.google.com",
  "username": "user@example.com",
  "password": "your_password",
  "category": "邮箱",
  "notes": "主要工作邮箱",
  "is_fav": false
}
```

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "created_at": "2025-10-31T14:00:00Z",
    "updated_at": "2025-10-31T14:00:00Z",
    "title": "Google账户",
    "url": "https://accounts.google.com",
    "username": "user@example.com",
    "category": "邮箱",
    "notes": "主要工作邮箱",
    "is_fav": false
  }
}
```

### 2. 获取账户列表

**GET** `/api/passwords?page=1&page_size=20&category=邮箱&keyword=google`

#### 查询参数
- `page`: 页码（默认：1）
- `page_size`: 每页条数（默认：20，最大：100）
- `category`: 分类筛选（可选）
- `keyword`: 关键词搜索（可选，搜索标题、URL、用户名、备注）

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 42,
    "items": [
      {
        "id": 1,
        "created_at": "2025-10-31T14:00:00Z",
        "updated_at": "2025-10-31T14:00:00Z",
        "title": "Google账户",
        "url": "https://accounts.google.com",
        "username": "user@example.com",
        "category": "邮箱",
        "notes": "主要工作邮箱",
        "is_fav": false
      }
    ]
  }
}
```

### 3. 获取账户详情

**GET** `/api/passwords/:id`

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "created_at": "2025-10-31T14:00:00Z",
    "updated_at": "2025-10-31T14:00:00Z",
    "title": "Google账户",
    "url": "https://accounts.google.com",
    "username": "user@example.com",
    "category": "邮箱",
    "notes": "主要工作邮箱",
    "is_fav": false
  }
}
```

**注意：** 密码不会在此接口返回，需要单独调用查看密码接口。

### 4. 查看密码

**GET** `/api/passwords/:id/password`

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "password": "your_password"
  }
}
```

**安全提示：** 此接口会记录日志，包括用户ID和账户ID。

### 5. 更新账户

**PUT** `/api/passwords/:id`

#### 请求体
```json
{
  "title": "Google账户（更新）",
  "url": "https://accounts.google.com",
  "username": "new_user@example.com",
  "password": "new_password",
  "category": "工作",
  "notes": "更新后的备注",
  "is_fav": true
}
```

**注意：**
- 所有字段都是可选的
- 如果 `password` 为空或不提供，则不更新密码
- `is_fav` 需要使用指针类型以区分未设置和 false

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "created_at": "2025-10-31T14:00:00Z",
    "updated_at": "2025-10-31T14:30:00Z",
    "title": "Google账户（更新）",
    "url": "https://accounts.google.com",
    "username": "new_user@example.com",
    "category": "工作",
    "notes": "更新后的备注",
    "is_fav": true
  }
}
```

### 6. 删除账户

**DELETE** `/api/passwords/:id`

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "删除成功"
  }
}
```

### 7. 获取分类列表

**GET** `/api/passwords/categories`

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "categories": [
      "邮箱",
      "社交",
      "工作",
      "购物",
      "金融"
    ]
  }
}
```

### 8. 获取统计信息

**GET** `/api/passwords/stats`

#### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 42,
    "fav_count": 5,
    "category_counts": {
      "邮箱": 8,
      "社交": 12,
      "工作": 15,
      "购物": 5,
      "金融": 2
    }
  }
}
```

## 📝 数据模型

### Account（账户）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键ID |
| created_at | time.Time | 创建时间 |
| updated_at | time.Time | 更新时间 |
| user_id | uint | 所属用户ID |
| title | string | 网站/应用标题（必填） |
| url | string | 网站URL |
| username | string | 用户名/账号 |
| password | string | 密码（加密存储，不返回） |
| category | string | 分类 |
| notes | string | 备注 |
| is_fav | bool | 是否收藏 |

## 🔒 加密说明

### 加密算法
- **算法**: AES-256-GCM
- **模式**: Galois/Counter Mode（提供认证加密）
- **密钥长度**: 256位（32字节）
- **输出格式**: Base64编码

### 环境变量

可以通过环境变量 `ENCRYPTION_KEY` 设置自定义加密密钥：

```bash
export ENCRYPTION_KEY="your-32-byte-encryption-key123"
```

**生产环境必须设置自己的密钥！**

### 加密流程

1. 使用 AES-256 初始化加密器
2. 生成随机 Nonce（96位）
3. 使用 GCM 模式加密明文
4. 将 Nonce 和密文组合
5. Base64 编码存储

### 解密流程

1. Base64 解码密文
2. 分离 Nonce 和加密数据
3. 使用 GCM 模式解密
4. 返回明文密码

## 🎯 使用场景

### 常见分类建议
- 邮箱：Gmail、Outlook等
- 社交：Facebook、Twitter、Instagram等
- 工作：企业内部系统、协作工具等
- 购物：淘宝、京东、Amazon等
- 金融：银行、支付宝、PayPal等
- 开发：GitHub、GitLab、云服务等
- 学习：Coursera、Udemy等

### 最佳实践

1. **分类管理**
   - 使用统一的分类名称
   - 不要创建过多分类

2. **安全建议**
   - 定期更新重要账户密码
   - 不要在备注中存储敏感信息
   - 使用收藏功能标记重要账户

3. **搜索优化**
   - 在标题中包含关键词
   - 善用备注字段记录额外信息

## ⚠️ 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权（Token无效或过期） |
| 404 | 账户不存在 |
| 500 | 服务器内部错误 |

## 📊 性能建议

1. **列表分页**
   - 默认每页20条
   - 建议不超过50条

2. **搜索优化**
   - 使用精确的关键词
   - 结合分类筛选提高效率

3. **收藏功能**
   - 列表默认按收藏优先排序
   - 常用账户建议标记收藏

---

**最后更新**: 2025-10-31
**版本**: 1.0.0

