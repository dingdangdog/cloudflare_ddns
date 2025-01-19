package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// 配置结构体
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
	CLOUDFLARE CloudflareConfig `json:"CLOUDFLARE"`
	IP_API_URL string           `json:"IP_API_URL"`
	INTERVAL   int              `json:"INTERVAL"`
}

// 获取配置文件
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

// 获取当前公网IP
func getPublicIP(ipApiUrl string) (string, error) {
	resp, err := http.Get(ipApiUrl)
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

// 更新 Cloudflare DNS 记录
func updateDNS(config *Config, ip string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", config.CLOUDFLARE.CF_ZONE_ID, config.CLOUDFLARE.CF_RECORD_ID)

	// 设置请求头
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", config.CLOUDFLARE.CF_API_TOKEN),
		"Content-Type":  "application/json",
	}

	// 请求数据
	data := fmt.Sprintf(`{
		"type": "%s",
		"name": "%s",
		"content": "%s",
		"ttl": %d,
		"proxied": %t
	}`, config.CLOUDFLARE.DNS_TYPE, config.CLOUDFLARE.DNS_DOMAIN_NAME, ip, config.CLOUDFLARE.DNS_TTL, config.CLOUDFLARE.DNS_PROXIED)

	// 创建 HTTP 请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v \n", err)
	}
	defer resp.Body.Close()

	now := time.Now().Format("2006/01/02 15:04:05")
	// 处理响应
	if resp.StatusCode == http.StatusOK {
		log.Printf("%s: DNS record updated successfully to IP: %s \n", now, ip)
	} else {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s: failed to update DNS: %v - %s \n", now, resp.StatusCode, string(body))
	}

	return nil
}

// 循环执行主逻辑
func main() {
	// 加载配置
	config, err := loadConfig("config.json")
	if err != nil {
		log.Printf("Error loading config: %v \n", err)
	}

	// 设置时间间隔
	interval := time.Duration(config.INTERVAL) * time.Second

	for {
		// log.Printf("Fetching current public IP...")
		ip, err := getPublicIP(config.IP_API_URL)
		if err != nil {
			log.Printf("Error fetching public IP: %v \n", err)
			time.Sleep(interval)
			continue
		}
		oldip, _ := os.ReadFile("ip.last")
		// 与原IP相同，跳过处理
		if ip == string(oldip) {
			log.Printf("Old IP: %s \n", ip)
			time.Sleep(interval)
			continue
		}
		// 保存新IP到文件
		_ = os.WriteFile("ip.last", []byte(ip), 0644)

		// log.Printf("Current public IP: %s \n", ip)

		// 更新 Cloudflare DNS 记录
		err = updateDNS(config, ip)
		if err != nil {
			log.Printf("Error updating DNS: %v \n", err)
		}

		// 等待下一次执行
		// log.Printf("Sleeping for %s \n", interval)
		time.Sleep(interval)
	}
}
