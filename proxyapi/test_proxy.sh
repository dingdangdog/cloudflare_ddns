#!/bin/bash

# Cloudflare DNS Proxy Service 测试脚本

# 配置变量
PROXY_URL="http://localhost:12322"  # 修改为你的代理服务地址
CLIENT_ID="0"
CLIENT_KEY="Test_Client_Key"

# 测试数据（请替换为真实数据）
CF_API_TOKEN="your_cf_api_token"
ZONE_ID="your_zone_id"
RECORD_ID="your_record_id"
DOMAIN_NAME="test.example.com"
IP_ADDRESS="1.2.3.4"

echo "=== Cloudflare DNS Proxy Service 测试 ==="
echo "代理服务地址: $PROXY_URL"
echo ""

# 1. 测试健康检查
echo "1. 测试健康检查..."
curl -s "$PROXY_URL/health" | jq .
echo ""

# 2. 测试默认端点
echo "2. 测试默认端点..."
curl -s "$PROXY_URL/" | jq .
echo ""

# 3. 测试DNS更新（需要真实数据）
echo "3. 测试DNS更新..."
echo "注意：请先修改脚本中的测试数据为真实值"

if [ "$CF_API_TOKEN" = "your_cf_api_token" ]; then
    echo "跳过DNS更新测试（需要配置真实数据）"
else
    curl -X POST "$PROXY_URL/update-dns?client_id=$CLIENT_ID&client_key=$CLIENT_KEY" \
        -H "Content-Type: application/json" \
        -d "{
            \"api_token\": \"$CF_API_TOKEN\",
            \"zone_id\": \"$ZONE_ID\",
            \"record_id\": \"$RECORD_ID\",
            \"type\": \"A\",
            \"name\": \"$DOMAIN_NAME\",
            \"content\": \"$IP_ADDRESS\",
            \"ttl\": 1,
            \"proxied\": false
        }" | jq .
fi

echo ""
echo "=== 测试完成 ==="
