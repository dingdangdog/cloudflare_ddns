package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// 处理请求并返回客户端的IP地址
func whoisme(w http.ResponseWriter, r *http.Request) {
	// 先从请求头获取 X-Forwarded-For 字段，如果没有则获取 RemoteAddr
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}
	// 如果IP后面有端口，去掉端口部分
	if strings.Contains(clientIP, ":") {
		clientIP = strings.Split(clientIP, ":")[0]
	}

	// 保存IP到文件
	err := os.WriteFile("ip.last", []byte("Client IP: "+clientIP), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving IP to file: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回请求者的IP
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(clientIP))
}

func lastip(w http.ResponseWriter, r *http.Request) {
	// 先从请求头获取 X-Forwarded-For 字段，如果没有则获取 RemoteAddr
	ip, err := os.ReadFile("ip.last")
	if err != nil {
		log.Println("error")
		return
	}
	// 返回请求者的IP
	w.WriteHeader(http.StatusOK)
	w.Write(ip)
}

func main() {
	http.HandleFunc("/whoisme", whoisme)
	http.HandleFunc("/lastip", lastip)
	log.Println("Server starting on port 12321...")
	err := http.ListenAndServe(":12321", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
