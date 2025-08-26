# Cloudflare DNS Proxy Service

这是一个用于解决客户端无法直接调用Cloudflare API问题的代理服务。

## 问题背景

在某些网络环境下，客户端可能无法直接访问Cloudflare API（`api.cloudflare.com`），出现SSL连接错误：
```
curl: (35) OpenSSL SSL_connect: SSL_ERROR_SYSCALL in connection to api.cloudflare.com:443
```

## 解决方案

通过代理服务中转，客户端不再直接调用Cloudflare API，而是：
1. 客户端 → 代理服务
2. 代理服务 → Cloudflare API
3. 代理服务 → 客户端

## 部署方式

### 1. Cloudflare Workers（推荐）

**优点**：
- 免费额度充足
- 全球CDN加速
- 无需服务器维护
- 自动HTTPS

**部署步骤**：
1. 在Cloudflare控制台创建Worker
2. 复制 `proxy.js` 代码到Worker
3. 配置环境变量 `CLIENT_KEYS`（多个密钥用逗号分隔）
4. 部署Worker

**使用示例**：
```bash
# 测试健康检查
curl https://your-worker.your-subdomain.workers.dev/health

# 更新DNS记录
curl -X POST "https://your-worker.your-subdomain.workers.dev/update-dns?client_id=0&client_key=your_key" \
  -H "Content-Type: application/json" \
  -d '{
    "api_token": "your_cf_token",
    "zone_id": "your_zone_id",
    "record_id": "your_record_id",
    "type": "A",
    "name": "example.com",
    "content": "1.2.3.4",
    "ttl": 1,
    "proxied": false
  }'
```

### 2. Go版本（Docker部署）

**优点**：
- 完全控制
- 可以部署在内网
- 支持更多自定义功能

**部署步骤**：
1. 复制 `config_demo.json` 为 `config.json` 并修改配置
2. 使用Docker部署：
   ```bash
   docker-compose up -d
   ```

**使用示例**：
```bash
# 测试健康检查
curl http://localhost:12322/health

# 更新DNS记录
curl -X POST "http://localhost:12322/update-dns?client_id=0&client_key=your_key" \
  -H "Content-Type: application/json" \
  -d '{
    "api_token": "your_cf_token",
    "zone_id": "your_zone_id",
    "record_id": "your_record_id",
    "type": "A",
    "name": "example.com",
    "content": "1.2.3.4",
    "ttl": 1,
    "proxied": false
  }'
```

## 配置说明

### 代理服务配置

**Cloudflare Workers**：
- 环境变量 `CLIENT_KEYS`：客户端密钥列表，用逗号分隔

**Go版本**：
```json
{
  "CLIENTS": ["key1", "key2", "key3"]
}
```

### DDNS客户端配置

在 `ddns/config.json` 中添加代理服务配置：
```json
{
  "PROXY_API_URL": "https://your-worker.your-subdomain.workers.dev",
  "PROXY_CLIENT_ID": 0,
  "PROXY_CLIENT_KEY": "your_key"
}
```

## API接口

### 1. 健康检查
- **路径**：`/health`
- **方法**：GET
- **响应**：
```json
{
  "status": "ok",
  "service": "cloudflare-dns-proxy",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 2. DNS更新
- **路径**：`/update-dns`
- **方法**：POST
- **参数**：
  - `client_id`：客户端ID（查询参数）
  - `client_key`：客户端密钥（查询参数）
- **请求体**：
```json
{
  "api_token": "your_cf_token",
  "zone_id": "your_zone_id",
  "record_id": "your_record_id",
  "type": "A",
  "name": "example.com",
  "content": "1.2.3.4",
  "ttl": 1,
  "proxied": false
}
```
- **响应**：
```json
{
  "success": true,
  "message": "DNS record updated successfully for example.com",
  "data": { ... }
}
```

## 安全说明

1. **客户端认证**：通过 `client_id` 和 `client_key` 进行认证
2. **HTTPS传输**：Cloudflare Workers自动提供HTTPS
3. **API令牌安全**：API令牌通过HTTPS传输，不会被明文记录

## 故障排除

1. **SSL错误**：确保代理服务可以访问 `api.cloudflare.com`
2. **认证失败**：检查 `client_id` 和 `client_key` 是否正确
3. **DNS更新失败**：检查Cloudflare API令牌权限和DNS记录ID是否正确
