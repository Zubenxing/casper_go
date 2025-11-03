# SSL 证书和私钥配对验证功能

## 📚 概述

SSL 证书验证功能现已增强，支持**证书和私钥的配对验证**，确保证书和私钥是匹配的，可以正常使用。

## 🎯 功能特性

### 1. 单独验证证书
只上传证书文件，验证：
- ✅ 证书格式是否正确
- ✅ 证书是否过期
- ✅ 证书链信息
- ✅ 证书用途和扩展信息

### 2. 证书-私钥配对验证 ⭐新增
同时上传证书和私钥，额外验证：
- ✅ 证书的公钥和私钥是否配对
- ✅ 证书是否可以和私钥一起使用
- ✅ 支持 RSA 和 ECDSA 算法

## 🔧 后端 API

### 请求示例

**仅验证证书：**
```bash
curl -X POST http://localhost:8080/api/certificates/tools/validate-cert \
  -H "Content-Type: application/json" \
  -d '{
    "cert_content": "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----"
  }'
```

**验证证书并校验私钥配对：**
```bash
curl -X POST http://localhost:8080/api/certificates/tools/validate-cert \
  -H "Content-Type: application/json" \
  -d '{
    "cert_content": "-----BEGIN CERTIFICATE-----\nMIID...\n-----END CERTIFICATE-----",
    "private_key_content": "-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----"
  }'
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "common_name": "example.com",
    "organization": "Example Corp",
    "dns_names": ["example.com", "*.example.com"],
    "not_before": "2024-01-01T00:00:00Z",
    "not_after": "2025-01-01T00:00:00Z",
    "days_left": 245,
    "is_expired": false,
    "key_pair_checked": true,
    "key_pair_matched": true,
    ...
  }
}
```

**响应字段说明：**
- `key_pair_checked` (bool): 是否进行了私钥配对检查
- `key_pair_matched` (bool): 证书和私钥是否配对成功

## 💻 前端使用

### 界面说明

在"证书监控" → "验证证书" 页面：

1. **证书文件区域**
   - 粘贴 PEM 格式的证书内容
   - 或拖拽 `.crt`、`.cer`、`.pem` 文件上传

2. **私钥文件区域（可选）** ⭐新增
   - 粘贴 PEM 格式的私钥内容
   - 支持 RSA、ECDSA 私钥
   - 提供后将自动进行配对验证

3. **验证结果**
   - 显示证书详细信息
   - 如果上传了私钥，会显示配对验证结果：
     - ✅ 绿色 Alert：证书和私钥配对成功
     - ❌ 红色 Alert：证书和私钥不匹配

### 配对验证逻辑

后端使用以下方法验证配对：

**RSA 密钥对：**
- 比较证书公钥的模数 (N) 和指数 (E)
- 与私钥的模数和指数是否相同

**ECDSA 密钥对：**
- 比较椭圆曲线类型
- 比较公钥坐标 (X, Y) 是否匹配

## 🧪 测试方法

### 1. 使用自己生成的证书和私钥

```bash
# 步骤 1: 使用"生成 CSR"功能生成 CSR 和私钥
# 步骤 2: 将 CSR 提交给 CA 获取证书
# 步骤 3: 使用"验证证书"功能，同时上传证书和私钥
```

### 2. 使用现有的证书和私钥

```bash
# 读取证书和私钥文件
CERT_CONTENT=$(cat /path/to/certificate.crt)
KEY_CONTENT=$(cat /path/to/private.key)

# 发送验证请求
curl -X POST http://localhost:8080/api/certificates/tools/validate-cert \
  -H "Content-Type: application/json" \
  -d "{
    \"cert_content\": \"$CERT_CONTENT\",
    \"private_key_content\": \"$KEY_CONTENT\"
  }"
```

### 3. 测试不匹配的私钥

为了测试验证功能，可以故意使用不匹配的私钥：

```bash
# 使用证书 A 和私钥 B（它们不配对）
# 应该返回 key_pair_matched: false
```

## 🔐 安全说明

1. **私钥安全**
   - 私钥内容仅用于验证，不会被保存
   - 验证后立即从内存中清除
   - 建议在安全的环境中使用此功能

2. **使用场景**
   - ✅ 配置 HTTPS 服务器前验证证书和私钥
   - ✅ 故障排查：检查证书配置错误
   - ✅ 证书更新：确认新证书和现有私钥匹配
   - ❌ 不要在公共网络上传私钥

## 📊 验证结果解读

| 场景 | key_pair_checked | key_pair_matched | 说明 |
|------|------------------|------------------|------|
| 只上传证书 | false | false | 仅验证证书本身 |
| 证书+私钥（匹配） | true | true | ✅ 证书可用 |
| 证书+私钥（不匹配） | true | false | ❌ 证书不可用 |

## 🎯 最佳实践

1. **部署前验证**
   ```bash
   # 在部署 SSL 证书到服务器前，先验证配对
   ```

2. **定期检查**
   ```bash
   # 定期检查现有证书和私钥的配对关系
   ```

3. **故障排查**
   ```bash
   # 当 HTTPS 服务器出现问题时，验证证书和私钥是否匹配
   ```

## 🐛 常见问题

### Q: 为什么私钥验证失败？
**A:** 可能原因：
1. 证书和私钥确实不配对
2. 私钥格式不正确（应为 PEM 格式）
3. 私钥加密了（需要先解密）

### Q: 支持哪些私钥格式？
**A:** 支持以下 PEM 格式：
- RSA PRIVATE KEY (PKCS#1)
- PRIVATE KEY (PKCS#8)
- EC PRIVATE KEY (ECDSA)

### Q: 私钥会被保存吗？
**A:** 不会。私钥仅在内存中用于验证，验证完成后立即释放。

## 📝 更新日志

- **2025-10-29**: 新增证书-私钥配对验证功能
- **2025-10-29**: 前端界面增加私钥上传区域
- **2025-10-29**: 后端增加 `verifyKeyPair` 函数

---

**提示**：这个功能大大提升了证书验证的实用性！现在不仅能验证证书本身，还能确保证书和私钥是可以一起使用的。🎉

