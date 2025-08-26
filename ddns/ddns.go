package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type CloudflareConfig struct {
	CF_API_TOKEN       string `json:"CF_API_TOKEN"`
	CF_ZONE_ID         string `json:"CF_ZONE_ID"`
	CF_RECORD_ID       string `json:"CF_RECORD_ID"`
	DNS_TYPE           string `json:"DNS_TYPE"`
	DNS_DOMAIN_NAME    string `json:"DNS_DOMAIN_NAME"`
	DNS_DOMAIN_CONTENT string `json:"DNS_DOMAIN_CONTENT"`
	DNS_TTL            int    `json:"DNS_TTL"`
	DNS_PROXIED        bool   `json:"DNS_PROXIED"`
}

type Config struct {
	CLOUDFLARE        []CloudflareConfig `json:"CLOUDFLARE"`
	WHOIAM_API_URL    string             `json:"WHOIAM_API_URL"`
	WHOIAM_CLIENT_ID  int                `json:"WHOIAM_CLIENT_ID"`
	WHOIAM_CLIENT_KEY string             `json:"WHOIAM_CLIENT_KEY"`
	PROXY_API_URL     string             `json:"PROXY_API_URL"`
	PROXY_CLIENT_ID   int                `json:"PROXY_CLIENT_ID"`
	PROXY_CLIENT_KEY  string             `json:"PROXY_CLIENT_KEY"`
	MODE              string             `json:"MODE"`
	INTERVAL          int                `json:"INTERVAL"`
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

// Get Public IP
func getPublicIP(config *Config) (string, error) {
	// 设置http get 请求的 url传参，参数为id 和 key
	url := fmt.Sprintf("%s?id=%d&key=%s", config.WHOIAM_API_URL, config.WHOIAM_CLIENT_ID, config.WHOIAM_CLIENT_KEY)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error getting public IP: %v", err)
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(ip), nil
}

// Update Cloudflare DNS records via proxy
func updateDNS(config *Config, ip string) error {
	for _, cfConfig := range config.CLOUDFLARE {
		domainNames := strings.Split(cfConfig.DNS_DOMAIN_NAME, ",")
		for _, domainName := range domainNames {
			domainName = strings.TrimSpace(domainName) // trim whitespace from domain name
			if domainName == "" {
				continue // skip empty domain names
			}

			// 构建代理请求URL
			proxyURL := fmt.Sprintf("%s/update-dns?client_id=%d&client_key=%s",
				config.PROXY_API_URL, config.PROXY_CLIENT_ID, config.PROXY_CLIENT_KEY)

			// 构建DNS更新请求体
			requestBody := map[string]interface{}{
				"api_token": cfConfig.CF_API_TOKEN,
				"zone_id":   cfConfig.CF_ZONE_ID,
				"record_id": cfConfig.CF_RECORD_ID,
				"type":      cfConfig.DNS_TYPE,
				"name":      domainName,
				"content":   ip,
				"ttl":       cfConfig.DNS_TTL,
				"proxied":   cfConfig.DNS_PROXIED,
			}

			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return fmt.Errorf("error marshaling JSON for domain %s: %v", domainName, err)
			}

			// 创建HTTP请求
			req, err := http.NewRequest("POST", proxyURL, bytes.NewBuffer(jsonData))
			if err != nil {
				return fmt.Errorf("error creating request for domain %s: %v", domainName, err)
			}

			// 设置请求头
			req.Header.Set("Content-Type", "application/json")

			// 发送请求
			client := &http.Client{
				Timeout: 30 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error sending request for domain %s: %v \n", domainName, err)
			}
			defer resp.Body.Close()

			// 读取响应
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response for domain %s: %v \n", domainName, err)
			}

			if resp.StatusCode == http.StatusOK {
				log.Printf("DNS record for domain %s updated successfully to IP: %s \n", domainName, ip)
			} else {
				log.Printf("Failed to update DNS for domain %s: Status %d, Response: %s \n", domainName, resp.StatusCode, string(body))
				return fmt.Errorf("failed to update DNS for domain %s: %v - %s \n", domainName, resp.StatusCode, string(body))
			}
		}
	}

	return nil
}

// Main
func main() {
	log.Printf("Starting \n")
	// load config
	config, err := loadConfig("config.json")
	if err != nil {
		log.Printf("Error loading config: %v \n", err)
	}
	// set sleep time
	interval := time.Duration(config.INTERVAL) * time.Second

	log.Printf("Started \n")
	for {
		// log.Printf("Fetching current public IP...")
		ip, err := getPublicIP(config)
		if err != nil {
			log.Printf("Error fetching public IP: %v \n", err)
			// 获取公网IP失败，立即重新执行
			continue
		}
		// 按,分割ip，保留第一个元素，并去除空格
		ip = strings.Split(ip, ",")[0]
		ip = strings.TrimSpace(ip)

		oldip, _ := os.ReadFile("ip.last")
		// save as the old IP, skip processing
		if ip == string(oldip) {
			// log.Printf("Old IP: %s \n", ip)
			// sleep
			time.Sleep(interval)
			continue
		}
		log.Printf("new IP: %s \n", ip)
		// log.Printf("Current public IP: %s \n", ip)

		// Update Cloudflare DNS records
		if config.MODE != "development" {
			err = updateDNS(config, ip)
			if err != nil {
				log.Printf("Error updating DNS: %v \n", err)
				// 出现任何错误，重新执行全部任务
				continue
			}

			// save the new IP to ip.last File
			_ = os.WriteFile("ip.last", []byte(ip), 0644)
			log.Printf("write new IP: %s \n", ip)
		}

		// sleep
		// log.Printf("Sleeping for %s \n", interval)
		time.Sleep(interval)
	}
}
