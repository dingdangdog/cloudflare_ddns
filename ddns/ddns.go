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
	CLOUDFLARE     CloudflareConfig `json:"CLOUDFLARE"`
	IP_API_URL     string           `json:"IP_API_URL"`
	DDD_CLIENT_ID  int              `json:"DDD_CLIENT_ID"`
	DDD_CLIENT_KEY string           `json:"DDD_CLIENT_KEY"`
	MODE           string           `json:"MODE"`
	INTERVAL       int              `json:"INTERVAL"`
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
	url := fmt.Sprintf("%s?id=%d&key=%s", config.IP_API_URL, config.DDD_CLIENT_ID, config.DDD_CLIENT_KEY)
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

// Update Cloudflare DNS records
func updateDNS(config *Config, ip string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", config.CLOUDFLARE.CF_ZONE_ID, config.CLOUDFLARE.CF_RECORD_ID)

	// request header
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", config.CLOUDFLARE.CF_API_TOKEN),
		"Content-Type":  "application/json",
	}

	// request body
	data := fmt.Sprintf(`{
		"type": "%s",
		"name": "%s",
		"content": "%s",
		"ttl": %d,
		"proxied": %t
	}`, config.CLOUDFLARE.DNS_TYPE, config.CLOUDFLARE.DNS_DOMAIN_NAME, ip, config.CLOUDFLARE.DNS_TTL, config.CLOUDFLARE.DNS_PROXIED)

	// create HTTP request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// set request header
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v \n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("DNS record updated successfully to IP: %s \n", ip)
	} else {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update DNS: %v - %s \n", resp.StatusCode, string(body))
	}

	return nil
}

// Main
func main() {
	// load config
	config, err := loadConfig("config.json")
	if err != nil {
		log.Printf("Error loading config: %v \n", err)
	}

	// set sleep time
	interval := time.Duration(config.INTERVAL) * time.Second

	for {
		// log.Printf("Fetching current public IP...")
		ip, err := getPublicIP(config)
		if err != nil {
			log.Printf("Error fetching public IP: %v \n", err)
			// sleep
			time.Sleep(interval)
			continue
		}
		// 按,分割ip，保留第一个元素，并去除空格
		ip = strings.Split(ip, ",")[0]
		ip = strings.TrimSpace(ip)

		oldip, _ := os.ReadFile("ip.last")
		// save as the old IP, skip processing
		if ip == string(oldip) {
			log.Printf("Old IP: %s \n", ip)
			// sleep
			time.Sleep(interval)
			continue
		}
		// save the new IP to ip.last File
		_ = os.WriteFile("ip.last", []byte(ip), 0644)

		// log.Printf("Current public IP: %s \n", ip)

		// Update Cloudflare DNS records
		if config.MODE != "development" {
			err = updateDNS(config, ip)
			if err != nil {
				log.Printf("Error updating DNS: %v \n", err)
			}
		}

		// sleep
		// log.Printf("Sleeping for %s \n", interval)
		time.Sleep(interval)
	}
}
