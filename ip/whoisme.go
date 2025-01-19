package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Processes the request and returns the client's IP address
func whoisme(w http.ResponseWriter, r *http.Request) {
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
	err := os.WriteFile("ip.last", []byte("Client IP: "+clientIP), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving IP to file: %v", err), http.StatusInternalServerError)
		return
	}

	// return the requester's IP
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(clientIP))
}

func lastip(w http.ResponseWriter, r *http.Request) {
	ip, err := os.ReadFile("ip.last")
	if err != nil {
		log.Println("error")
		return
	}
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
