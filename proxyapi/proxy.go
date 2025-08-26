package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	CLIENTS []string `json:"CLIENTS"`
}

type DNSUpdateRequest struct {
	APIToken string `json:"api_token"`
	ZoneID   string `json:"zone_id"`
	RecordID string `json:"record_id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Proxied  bool   `json:"proxied"`
}

type DNSUpdateResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Get json file content
func loadConfig(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}

	return &config, nil
}

// 验证客户端
func validateClient(r *http.Request, config *Config) (bool, string) {
	clientID := r.URL.Query().Get("client_id")
	clientKey := r.URL.Query().Get("client_key")

	if clientID == "" || clientKey == "" {
		return false, "Missing client_id or client_key"
	}

	idInt, err := strconv.Atoi(clientID)
	if err != nil {
		return false, "Invalid client_id"
	}

	if idInt < 0 || idInt >= len(config.CLIENTS) {
		return false, "Invalid client_id"
	}

	if config.CLIENTS[idInt] != clientKey {
		return false, "Invalid client_key"
	}

	return true, ""
}

// 调用Cloudflare API更新DNS记录
func updateCloudflareDNS(dnsData *DNSUpdateRequest) (*DNSUpdateResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", dnsData.ZoneID, dnsData.RecordID)

	// 构建请求体
	requestBody := map[string]interface{}{
		"type":    dnsData.Type,
		"name":    dnsData.Name,
		"content": dnsData.Content,
		"ttl":     dnsData.TTL,
		"proxied": dnsData.Proxied,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dnsData.APIToken))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var response DNSUpdateResponse
	if resp.StatusCode == http.StatusOK {
		response.Success = true
		response.Message = fmt.Sprintf("DNS record updated successfully for %s", dnsData.Name)

		// 解析Cloudflare API响应
		var cfResponse map[string]interface{}
		if err := json.Unmarshal(body, &cfResponse); err == nil {
			response.Data = cfResponse
		}
	} else {
		response.Success = false
		response.Message = fmt.Sprintf("Failed to update DNS record for %s", dnsData.Name)

		// 解析错误响应
		var cfError map[string]interface{}
		if err := json.Unmarshal(body, &cfError); err == nil {
			response.Error = cfError
		} else {
			response.Error = string(body)
		}
	}

	return &response, nil
}

// 处理DNS更新请求
func handleDNSUpdate(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 只允许POST方法
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 加载配置
	config, err := loadConfig("config.json")
	if err != nil {
		log.Printf("Error loading config: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 验证客户端
	valid, errorMsg := validateClient(r, config)
	if !valid {
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	// 解析请求体
	var dnsRequest DNSUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&dnsRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 验证必需字段
	if dnsRequest.APIToken == "" || dnsRequest.ZoneID == "" || dnsRequest.RecordID == "" ||
		dnsRequest.Type == "" || dnsRequest.Name == "" || dnsRequest.Content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 设置默认值
	if dnsRequest.TTL == 0 {
		dnsRequest.TTL = 1
	}

	// 调用Cloudflare API
	response, err := updateCloudflareDNS(&dnsRequest)
	if err != nil {
		log.Printf("Error updating DNS: %v", err)
		errorResponse := DNSUpdateResponse{
			Success: false,
			Message: "Internal server error",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

// 健康检查端点
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"status":    "ok",
		"service":   "cloudflare-dns-proxy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// 默认端点
func handleDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"service": "cloudflare-dns-proxy",
		"endpoints": map[string]string{
			"health":     "/health",
			"update_dns": "/update-dns",
		},
		"usage": "Send POST request to /update-dns with DNS data and client authentication",
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/update-dns", handleDNSUpdate)
	http.HandleFunc("/health", handleHealthCheck)
	http.HandleFunc("/", handleDefault)

	log.Println("Cloudflare DNS Proxy Server starting on port 12322...")
	err := http.ListenAndServe(":12322", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
