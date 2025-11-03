# 前端接口文档

**服务地址**: `http://localhost:8080`  
**API 前缀**: `/api`

## 💡 证书监控说明

本系统使用 **Go 原生 TLS 库**直接检查 SSL 证书，无需依赖外部组件。

**特点：**
- ✅ 独立运行，不依赖 Prometheus/blackbox_exporter
- ✅ 证书信息更详细（颁发者、组织、Subject等）
- ✅ 实时检查，响应快速
- ✅ 支持手动添加任意 HTTPS URL
- ✅ 后台定时自动更新所有证书状态

---

## 🔑 认证说明

除了登录、刷新令牌接口外，所有接口都需要在请求头中携带 Token：

```
Authorization: Bearer {你的token}
```

---

## 📋 接口列表

### 1. 用户认证

#### 1.1 用户登录
- **地址**: `POST /api/auth/login`
- **参数**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```
- **成功响应**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_at": "2024-10-28T15:35:38+08:00",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@casper.local",
      "nickname": "系统管理员",
      "role": "admin",
      "status": 1,
      "created_at": "2024-10-27T15:35:38+08:00"
    }
  }
}
```

#### 1.2 刷新令牌
- **地址**: `POST /api/auth/refresh`
- **参数**:
```json
{
  "refresh_token": "你的refresh_token"
}
```
- **成功响应**: 同登录接口

#### 1.3 用户登出
- **地址**: `POST /api/logout`
- **需要认证**: ✅
- **成功响应**:
```json
{
  "code": 0,
  "message": "登出成功",
  "data": null
}
```

#### 1.4 获取当前用户信息
- **地址**: `GET /api/profile`
- **需要认证**: ✅
- **成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@casper.local",
    "nickname": "系统管理员",
    "role": "admin",
    "status": 1,
    "created_at": "2024-10-27T15:35:38+08:00"
  }
}
```

---

### 2. 证书监控

#### 2.1 添加证书监控
- **地址**: `POST /api/certificates`
- **需要认证**: ✅
- **参数**:
```json
{
  "url": "https://www.baidu.com"
}
```
- **成功响应**:
```json
{
  "code": 0,
  "message": "添加成功",
  "data": {
    "id": 1,
    "url": "https://www.baidu.com",
    "domain": "www.baidu.com",
    "issuer": "GlobalSign RSA OV SSL CA 2018",
    "subject": "www.baidu.com",
    "organization": "Beijing Baidu Netcom Science Technology Co., Ltd",
    "not_before": "2024-05-01T00:00:00Z",
    "not_after": "2025-06-01T23:59:59Z",
    "days_left": 217,
    "is_valid": true,
    "status": 1,
    "status_text": "正常",
    "error_msg": "",
    "last_check_at": "2024-10-27T15:35:40+08:00",
    "created_at": "2024-10-27T15:35:40+08:00"
  }
}
```

#### 2.2 获取证书监控列表（分页）
- **地址**: `GET /api/certificates?page=1&page_size=10`
- **需要认证**: ✅
- **查询参数**:
  - `page`: 页码（默认1）
  - `page_size`: 每页数量（默认10，最大100）
- **成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "url": "https://www.baidu.com",
        "domain": "www.baidu.com",
        "issuer": "GlobalSign RSA OV SSL CA 2018",
        "organization": "Beijing Baidu Netcom Science Technology Co., Ltd",
        "not_after": "2025-06-01T23:59:59Z",
        "days_left": 217,
        "is_valid": true,
        "status": 1,
        "status_text": "正常",
        "last_check_at": "2024-10-27T15:35:40+08:00"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

#### 2.3 获取证书详情
- **地址**: `GET /api/certificates/:id`
- **需要认证**: ✅
- **成功响应**: 同添加证书监控的 data 部分

#### 2.4 更新证书信息（手动触发检查）
- **地址**: `PUT /api/certificates/:id`
- **需要认证**: ✅
- **成功响应**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "url": "https://www.baidu.com",
    "days_left": 217,
    "status": 1,
    "status_text": "正常"
    // ... 其他字段
  }
}
```

#### 2.5 删除证书监控
- **地址**: `DELETE /api/certificates/:id`
- **需要认证**: ✅
- **成功响应**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

#### 2.6 检查所有证书（手动触发）
- **地址**: `POST /api/certificates/check-all`
- **需要认证**: ✅
- **成功响应**:
```json
{
  "code": 0,
  "message": "检查完成",
  "data": null
}
```

---

### 3. 其他接口

#### 3.1 健康检查
- **地址**: `GET /health`
- **无需认证**
- **成功响应**:
```json
{
  "status": "ok",
  "message": "Casper Platform is running"
}
```

---

## 🔢 状态码说明

| 状态码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1000 | 请求参数错误 |
| 1001 | 未授权 |
| 1002 | 禁止访问 |
| 1003 | 资源不存在 |
| 1005 | 请求过于频繁（限流） |
| 2000 | 用户名或密码错误 |
| 2001 | 无效的令牌 |
| 2002 | 令牌已过期 |
| 2003 | 令牌已失效 |
| 2004 | 用户已被禁用 |
| 3000 | URL 已存在 |
| 3001 | 证书记录不存在 |

---

## 🎯 证书状态说明

| status | status_text | 说明 |
|--------|-------------|------|
| 1 | 正常 | 证书有效，剩余天数充足 |
| 2 | 警告 | 证书即将过期（默认30天内） |
| 3 | 过期 | 证书已过期 |
| 0 | 禁用 | 监控已禁用 |

---

## 📦 响应格式

所有接口统一返回格式：

```json
{
  "code": 0,           // 状态码，0 表示成功
  "message": "success",// 消息
  "data": {}           // 数据（可选）
}
```

---

## 🚀 快速测试

### JavaScript / Axios 示例

```javascript
// 1. 登录
const loginResponse = await axios.post('http://localhost:8080/api/auth/login', {
  username: 'admin',
  password: 'admin123'
});

const token = loginResponse.data.data.access_token;

// 2. 设置默认请求头
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

// 3. 获取证书列表
const certList = await axios.get('http://localhost:8080/api/certificates?page=1&page_size=10');
console.log(certList.data);

// 4. 添加证书监控
const addCert = await axios.post('http://localhost:8080/api/certificates', {
  url: 'https://www.baidu.com'
});
```

### Vue 3 示例

```javascript
// api.js
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
});

// 请求拦截器
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器
api.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      // Token 失效，跳转登录
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 导出接口
export default {
  // 登录
  login: (username, password) => 
    api.post('/auth/login', { username, password }),
  
  // 获取用户信息
  getProfile: () => 
    api.get('/profile'),
  
  // 证书列表
  getCertificates: (page = 1, pageSize = 10) => 
    api.get(`/certificates?page=${page}&page_size=${pageSize}`),
  
  // 添加证书
  addCertificate: (url) => 
    api.post('/certificates', { url }),
  
  // 删除证书
  deleteCertificate: (id) => 
    api.delete(`/certificates/${id}`)
};
```

---

## ⚠️ 注意事项

1. **Token 过期时间**: 24小时，过期后使用 refresh_token 刷新
2. **限流**: 每个IP每分钟最多100个请求
3. **请求ID**: 每个响应都会返回 `X-Request-ID` 响应头，用于追踪
4. **CORS**: 已配置允许跨域
5. **默认账号**: username: `admin`, password: `admin123`（请及时修改）

---

## 📖 完整文档

- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **项目文档**: 查看 README.md

---

**最后更新**: 2024-10-27

