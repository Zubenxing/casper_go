#!/bin/bash

# 测试证书和私钥配对验证
# 这个脚本演示如何验证证书和私钥是否匹配

echo "================================"
echo "测试 1: 使用配对的证书和私钥"
echo "================================"

# 首先生成一个证书和私钥对
CSR_RESULT=$(curl -s -X POST http://localhost:8080/api/certificates/tools/generate-csr \
  -H "Content-Type: application/json" \
  -d '{
    "common_name": "test.example.com",
    "country": "CN",
    "province": "Beijing",
    "locality": "Beijing",
    "organization": "Test Org",
    "organizational_unit": "IT",
    "key_algorithm": "RSA",
    "key_size": 2048
  }')

if [ $? -eq 0 ]; then
    echo "✓ CSR 生成成功"
    
    # 提取私钥（这里仅用于测试，实际生产中不会这样做）
    PRIVATE_KEY=$(echo "$CSR_RESULT" | jq -r '.data.private_key')
    CSR=$(echo "$CSR_RESULT" | jq -r '.data.csr')
    
    echo ""
    echo "注意：这个测试使用 CSR（证书签名请求），而不是真正的证书。"
    echo "在实际使用中，你需要用真实的证书和私钥来测试。"
    echo ""
    
    # 保存到临时文件
    echo "$PRIVATE_KEY" > /tmp/test_private_key.pem
    echo "$CSR" > /tmp/test_csr.pem
    
    echo "私钥已保存到: /tmp/test_private_key.pem"
    echo "CSR 已保存到: /tmp/test_csr.pem"
    echo ""
else
    echo "✗ CSR 生成失败"
    exit 1
fi

echo "================================"
echo "测试 2: 验证真实的证书（需要手动提供）"
echo "================================"
echo ""
echo "提示：请准备一对真实的证书和私钥文件进行测试。"
echo ""
echo "测试命令示例："
echo ""
echo "curl -X POST http://localhost:8080/api/certificates/tools/validate-cert \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d @- <<'EOF'"
echo "{"
echo "  \"cert_content\": \"$(cat your-certificate.crt)\","
echo "  \"private_key_content\": \"$(cat your-private-key.key)\""
echo "}"
echo "EOF"
echo ""
echo "================================"
echo "测试完成"
echo "================================"

