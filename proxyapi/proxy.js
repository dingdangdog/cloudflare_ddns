// Cloudflare DNS Proxy Service
// 用于解决客户端无法直接调用Cloudflare API的问题

// 验证客户端密钥
function validateClient(request, env) {
  const url = new URL(request.url);
  const clientId = url.searchParams.get("client_id");
  const clientKey = url.searchParams.get("client_key");

  // 从环境变量获取客户端配置
  const clients = env.CLIENT_KEYS ? env.CLIENT_KEYS.split(",") : [];

  // 验证client_id
  const idInt = parseInt(clientId, 10);
  if (isNaN(idInt) || idInt < 0 || idInt >= clients.length) {
    return { valid: false, error: "Invalid client_id" };
  }

  // 验证client_key
  if (clients[idInt] !== clientKey) {
    return { valid: false, error: "Invalid client_key" };
  }

  return { valid: true };
}

// 调用Cloudflare API更新DNS记录
async function updateCloudflareDNS(dnsData) {
  const url = `https://api.cloudflare.com/client/v4/zones/${dnsData.zone_id}/dns_records/${dnsData.record_id}`;

  const requestBody = {
    type: dnsData.type,
    name: dnsData.name,
    content: dnsData.content,
    ttl: dnsData.ttl,
    proxied: dnsData.proxied,
  };

  const response = await fetch(url, {
    method: "PATCH",
    headers: {
      Authorization: `Bearer ${dnsData.api_token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestBody),
  });

  const result = await response.json();

  return {
    status: response.status,
    success: response.ok,
    data: result,
  };
}

// 处理DNS更新请求
async function handleDNSUpdate(request, env) {
  try {
    // 验证客户端
    const validation = validateClient(request, env);
    if (!validation.valid) {
      return new Response(validation.error, { status: 400 });
    }

    // 解析请求体
    const dnsData = await request.json();

    // 验证必需字段
    const requiredFields = [
      "api_token",
      "zone_id",
      "record_id",
      "type",
      "name",
      "content",
    ];
    for (const field of requiredFields) {
      if (!dnsData[field]) {
        return new Response(`Missing required field: ${field}`, {
          status: 400,
        });
      }
    }

    // 设置默认值
    if (dnsData.ttl === undefined) dnsData.ttl = 1;
    if (dnsData.proxied === undefined) dnsData.proxied = false;

    // 调用Cloudflare API
    const result = await updateCloudflareDNS(dnsData);

    if (result.success) {
      return new Response(
        JSON.stringify({
          success: true,
          message: `DNS record updated successfully for ${dnsData.name}`,
          data: result.data,
        }),
        {
          status: 200,
          headers: { "Content-Type": "application/json" },
        }
      );
    } else {
      return new Response(
        JSON.stringify({
          success: false,
          message: `Failed to update DNS record for ${dnsData.name}`,
          error: result.data,
        }),
        {
          status: result.status,
          headers: { "Content-Type": "application/json" },
        }
      );
    }
  } catch (error) {
    return new Response(
      JSON.stringify({
        success: false,
        message: "Internal server error",
        error: error.message,
      }),
      {
        status: 500,
        headers: { "Content-Type": "application/json" },
      }
    );
  }
}

// 健康检查端点
function handleHealthCheck() {
  return new Response(
    JSON.stringify({
      status: "ok",
      service: "cloudflare-dns-proxy",
      timestamp: new Date().toISOString(),
    }),
    {
      status: 200,
      headers: { "Content-Type": "application/json" },
    }
  );
}

// 主处理函数
export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    // 健康检查
    if (url.pathname === "/health") {
      return handleHealthCheck();
    }

    // DNS更新端点
    if (url.pathname === "/update-dns") {
      if (request.method !== "POST") {
        return new Response("Method not allowed", { status: 405 });
      }
      return handleDNSUpdate(request, env);
    }

    // 默认端点（兼容性）
    if (url.pathname === "/") {
      return new Response(
        JSON.stringify({
          service: "cloudflare-dns-proxy",
          endpoints: {
            health: "/health",
            update_dns: "/update-dns",
          },
          usage:
            "Send POST request to /update-dns with DNS data and client authentication",
        }),
        {
          status: 200,
          headers: { "Content-Type": "application/json" },
        }
      );
    }

    return new Response("Not Found", { status: 404 });
  },
};
