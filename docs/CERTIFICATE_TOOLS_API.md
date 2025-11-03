# 证书工具 API 文档

## 📋 概述

证书工具提供三个核心功能的在线工具：
1. **CSR 生成** - 生成证书签名请求和私钥
2. **CSR 验证** - 验证 CSR 文件的合法性
3. **证书验证** - 验证 SSL 证书文件并提取详细信息

所有工具 API **不需要认证**，可公开访问。

---

## 🔧 API 端点

### 1. 生成 CSR

**接口**: `POST /api/certificates/tools/generate-csr`

**描述**: 在线生成 SSL 证书签名请求（CSR）和私钥

#### 请求参数

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| `common_name` | string | ✅ | 通用名/域名 | `example.com` |
| `country` | string | ❌ | 国家代码（2位） | `CN`, `US` |
| `province` | string | ❌ | 省份/州 | `Beijing`, `California` |
| `locality` | string | ❌ | 城市 | `Beijing`, `San Francisco` |
| `organization` | string | ❌ | 组织名称 | `Example Corp` |
| `organizational_unit` | string | ❌ | 部门 | `IT Department` |
| `email_address` | string | ❌ | 邮箱地址 | `admin@example.com` |
| `dns_names` | array | ❌ | SAN - DNS 名称列表 | `["example.com", "www.example.com"]` |
| `ip_addresses` | array | ❌ | SAN - IP 地址列表 | `["192.168.1.1"]` |
| `key_algorithm` | string | ✅ | 密钥算法 | `RSA` 或 `ECDSA` |
| `key_size` | int | ✅ | 密钥长度 | RSA: `2048`/`4096`, ECDSA: `256`/`384` |

#### 请求示例

**RSA 2048:**
```json
{
  "common_name": "example.com",
  "country": "CN",
  "province": "Beijing",
  "locality": "Beijing",
  "organization": "Example Corp",
  "organizational_unit": "IT Department",
  "email_address": "admin@example.com",
  "dns_names": ["example.com", "www.example.com"],
  "key_algorithm": "RSA",
  "key_size": 2048
}
```

**ECDSA P-256:**
```json
{
  "common_name": "ecdsa.example.com",
  "country": "US",
  "province": "California",
  "locality": "San Francisco",
  "organization": "ECDSA Test Corp",
  "key_algorithm": "ECDSA",
  "key_size": 256
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "CSR 生成成功",
  "data": {
    "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----\n",
    "csr": "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----\n"
  }
}
```

#### 响应字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `private_key` | string | PEM 格式的私钥 |
| `csr` | string | PEM 格式的 CSR |

#### cURL 示例

```bash
curl -X POST http://localhost:8080/api/certificates/tools/generate-csr \
  -H "Content-Type: application/json" \
  -d '{
    "common_name": "example.com",
    "country": "CN",
    "organization": "Example Corp",
    "key_algorithm": "RSA",
    "key_size": 2048
  }'
```

---

### 2. 验证 CSR

**接口**: `POST /api/certificates/tools/validate-csr`

**描述**: 验证 CSR 文件的合法性并提取详细信息

#### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `csr_content` | string | ✅ | PEM 格式的 CSR 内容 |

#### 请求示例

```json
{
  "csr_content": "-----BEGIN CERTIFICATE REQUEST-----\nMIIDCDCCAfACAQAwdjELMAkGA1UEBhMCQ04...\n-----END CERTIFICATE REQUEST-----"
}
```

#### 响应示例（成功）

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "common_name": "example.com",
    "country": "CN",
    "province": "Beijing",
    "locality": "Beijing",
    "organization": "Example Corp",
    "organizational_unit": "IT Department",
    "dns_names": ["example.com", "www.example.com"],
    "email_addresses": ["admin@example.com"],
    "signature_algorithm": "SHA256-RSA",
    "public_key_algorithm": "RSA",
    "public_key_size": 2048
  }
}
```

#### 响应示例（失败）

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": false,
    "error_message": "无效的 PEM 格式"
  }
}
```

#### 响应字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `valid` | boolean | CSR 是否有效 |
| `common_name` | string | 通用名 |
| `country` | string | 国家 |
| `province` | string | 省份 |
| `locality` | string | 城市 |
| `organization` | string | 组织 |
| `organizational_unit` | string | 部门 |
| `dns_names` | array | DNS 名称列表 |
| `email_addresses` | array | 邮箱地址列表 |
| `signature_algorithm` | string | 签名算法 |
| `public_key_algorithm` | string | 公钥算法 |
| `public_key_size` | int | 公钥长度（bit） |
| `error_message` | string | 错误信息（仅当 valid=false 时） |

#### cURL 示例

```bash
curl -X POST http://localhost:8080/api/certificates/tools/validate-csr \
  -H "Content-Type: application/json" \
  -d @test_data/validate_csr.json
```

---

### 3. 验证证书

**接口**: `POST /api/certificates/tools/validate-cert`

**描述**: 在线验证 SSL 证书文件并提取详细信息

#### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `cert_content` | string | ✅ | PEM 格式的证书内容 |

#### 请求示例

```json
{
  "cert_content": "-----BEGIN CERTIFICATE-----\nMIIDrzCCApegAwIBAgIQCDvgVpBCRrGhdWrJWZHHSjANBgkqhkiG9w0BAQUFADBh...\n-----END CERTIFICATE-----"
}
```

#### 响应示例（成功）

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "version": 3,
    "serial_number": "12345678901234567890",
    "issuer": "CN=DigiCert SHA2 Secure Server CA,O=DigiCert Inc,C=US",
    "subject": "CN=example.com,O=Example Organization,L=San Francisco,ST=California,C=US",
    "common_name": "example.com",
    "organization": "Example Organization",
    "country": "US",
    "dns_names": ["example.com", "www.example.com"],
    "ip_addresses": [],
    "email_addresses": [],
    "not_before": "2024-01-01T00:00:00Z",
    "not_after": "2025-01-01T23:59:59Z",
    "days_left": 365,
    "is_expired": false,
    "is_not_yet_valid": false,
    "signature_algorithm": "SHA256-RSA",
    "public_key_algorithm": "RSA",
    "public_key_size": 2048,
    "key_usage": ["Digital Signature", "Key Encipherment"],
    "ext_key_usage": ["Server Authentication", "Client Authentication"],
    "is_ca": false
  }
}
```

#### 响应示例（失败）

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": false,
    "error_message": "无效的 PEM 格式"
  }
}
```

#### 响应字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `valid` | boolean | 证书是否有效 |
| `version` | int | 证书版本 |
| `serial_number` | string | 序列号 |
| `issuer` | string | 颁发者 DN |
| `subject` | string | 主题 DN |
| `common_name` | string | 通用名 |
| `organization` | string | 组织 |
| `country` | string | 国家 |
| `dns_names` | array | SAN - DNS 名称 |
| `ip_addresses` | array | SAN - IP 地址 |
| `email_addresses` | array | SAN - 邮箱地址 |
| `not_before` | datetime | 生效时间 |
| `not_after` | datetime | 过期时间 |
| `days_left` | int | 剩余天数 |
| `is_expired` | boolean | 是否已过期 |
| `is_not_yet_valid` | boolean | 是否尚未生效 |
| `signature_algorithm` | string | 签名算法 |
| `public_key_algorithm` | string | 公钥算法 |
| `public_key_size` | int | 公钥长度 |
| `key_usage` | array | 密钥用途 |
| `ext_key_usage` | array | 扩展密钥用途 |
| `is_ca` | boolean | 是否为 CA 证书 |
| `error_message` | string | 错误信息（仅当 valid=false 时） |

---

## 🚀 使用场景

### 场景 1: 生成证书申请文件

1. 调用 `generate-csr` 接口生成 CSR 和私钥
2. 保存私钥到安全位置
3. 将 CSR 提交给证书颁发机构（CA）

### 场景 2: 验证 CSR 文件

1. 上传或粘贴 CSR 内容
2. 调用 `validate-csr` 接口验证
3. 查看 CSR 详细信息（主题、公钥等）

### 场景 3: 验证已有证书

1. 上传或粘贴证书内容
2. 调用 `validate-cert` 接口验证
3. 查看证书详情（颁发者、有效期、SAN 等）

---

## ⚠️ 注意事项

1. **私钥安全**: 生成的私钥仅在响应中返回一次，请务必妥善保存
2. **PEM 格式**: 所有证书和 CSR 必须是 PEM 格式（Base64 编码）
3. **密钥长度**:
   - RSA 推荐 2048 位或 4096 位
   - ECDSA 推荐 P-256 (256位) 或 P-384 (384位)
4. **国家代码**: 必须是 ISO 3166-1 alpha-2 标准的两位代码
5. **公开访问**: 这些 API 不需要认证，任何人都可以访问

---

## 📝 错误码

| Code | Message | 说明 |
|------|---------|------|
| 0 | success | 成功 |
| 400 | 请求参数错误 | 参数验证失败 |
| 500 | 服务器错误 | 内部处理错误 |

---

## 🧪 测试数据

测试文件位于 `test_data/` 目录：
- `generate_csr.json` - RSA CSR 生成测试数据
- `generate_csr_ecdsa.json` - ECDSA CSR 生成测试数据
- `validate_csr.json` - CSR 验证测试数据

---

## 🔗 相关资源

- [X.509 证书标准](https://tools.ietf.org/html/rfc5280)
- [PKCS#10 CSR 规范](https://tools.ietf.org/html/rfc2986)
- [PEM 格式说明](https://tools.ietf.org/html/rfc7468)

