package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	CLIENTS []string `json:"CLIENTS"`
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

// Processes the request and returns the client's IP address
func whoiam(w http.ResponseWriter, r *http.Request) {
	// load config
	configs, err := loadConfig("config.json")
	if err != nil {
		log.Printf("Error loading config: %v \n", err)
	}
	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	key := r.URL.Query().Get("key")
	if configs.CLIENTS[idInt] != key {
		http.Error(w, "Key Error", http.StatusBadRequest)
		return
	}

	// get the X-Forwarded-For field from the request header, if not, get the RemoteAddr
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}
	// if there is a port after the IP, remove the port part
	if strings.Contains(clientIP, ":") {
		clientIP = strings.Split(clientIP, ":")[0]
	}

	// save IP to ip.last File
	err = os.WriteFile(id+".ip", []byte("Client IP: "+clientIP), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving IP to file: %v", err), http.StatusInternalServerError)
		return
	}

	// return the requester's IP
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(clientIP))
}

func lastip(w http.ResponseWriter, r *http.Request) {
	// load config
	configs, err := loadConfig("config.json")
	if err != nil {
		log.Printf("Error loading config: %v \n", err)
	}
	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	key := r.URL.Query().Get("key")
	if configs.CLIENTS[idInt] != key {
		http.Error(w, "Key Error", http.StatusBadRequest)
		return
	}

	ip, err := os.ReadFile(id + ".ip")
	if err != nil {
		log.Println("error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(ip)
}

func main() {
	http.HandleFunc("/whoiam", whoiam)
	http.HandleFunc("/lastip", lastip)
	log.Println("Server starting on port 12321...")
	err := http.ListenAndServe(":12321", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
