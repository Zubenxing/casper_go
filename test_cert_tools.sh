#!/bin/bash

# 证书工具 API 测试脚本

BASE_URL="http://localhost:8080/api"

echo "================================"
echo "证书工具 API 测试"
echo "================================"
echo ""

# 1. 测试 CSR 生成
echo "1. 测试 CSR 生成..."
echo "---"
curl -X POST "${BASE_URL}/certificates/tools/generate-csr" \
  -H "Content-Type: application/json" \
  -d '{
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
  }' | jq .
echo ""
echo ""

# 2. 测试 CSR 验证（需要先生成一个 CSR）
echo "2. 测试 CSR 验证..."
echo "---"

# 生成一个测试 CSR
CSR_DATA=$(curl -s -X POST "${BASE_URL}/certificates/tools/generate-csr" \
  -H "Content-Type: application/json" \
  -d '{
    "common_name": "test.example.com",
    "country": "US",
    "province": "California",
    "locality": "San Francisco",
    "organization": "Test Org",
    "key_algorithm": "RSA",
    "key_size": 2048
  }')

# 提取 CSR 内容
CSR_CONTENT=$(echo "$CSR_DATA" | jq -r '.data.csr')

# 验证 CSR
curl -X POST "${BASE_URL}/certificates/tools/validate-csr" \
  -H "Content-Type: application/json" \
  -d "{
    \"csr_content\": $(echo "$CSR_CONTENT" | jq -Rs .)
  }" | jq .
echo ""
echo ""

# 3. 测试证书验证（使用一个示例证书）
echo "3. 测试证书验证..."
echo "---"

# 这里需要一个真实的 PEM 格式证书进行测试
# 可以使用 openssl s_client -connect example.com:443 -servername example.com | openssl x509 -text
echo "注意：此测试需要真实的证书内容，跳过..."
echo ""

echo "================================"
echo "测试完成！"
echo "================================"

