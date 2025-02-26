// Get client IP
function getClientIP(request) {
  let clientIP =
    request.headers.get("X-Forwarded-For") ||
    request.headers.get("CF-Connecting-IP");
  if (!clientIP) {
    clientIP = "unknown"; // If no IP information is available, return "unknown"
  }
  // If there are multiple IPs (X-Forwarded-For may contain a proxy chain), take the first one
  if (clientIP.includes(",")) {
    clientIP = clientIP.split(",")[0].trim();
  }
  return clientIP;
}

// Handle the /whoiam endpoint
async function handleWhoiam(request, env) {
  const url = new URL(request.url);
  const id = url.searchParams.get("id");
  const key = url.searchParams.get("key");

  // Get CLIENTS from environment variable (convert comma-separated string to array)
  const clients = env.CLIENT_KEYS ? env.CLIENT_KEYS.split(",") : [];

  // Validate id
  const idInt = parseInt(id, 10);
  if (isNaN(idInt) || idInt < 0 || idInt >= clients.length) {
    return new Response("Invalid ID", { status: 400 });
  }

  // Validate key
  if (clients[idInt] !== key) {
    return new Response("Key Error", { status: 400 });
  }

  // Get client IP and return it
  const clientIP = getClientIP(request);
  return new Response(clientIP, { status: 200 });
}

// Main handler function
export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname === "/") {
      return handleWhoiam(request, env);
    } else {
      return new Response("Not Found", { status: 404 });
    }
  },
};
