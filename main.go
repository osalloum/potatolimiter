package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	maxConcurrentRequests = 10
	queueTimeout          = 30 * time.Second
	upstreamURL           = "https://2ypozmzsat3qefvg4aacdn2vxu0bhgfl.lambda-url.eu-west-1.on.aws" // Replace with your upstream server URL
)

var (
	// Semaphore to limit concurrent requests
	concurrentRequests = make(chan struct{}, maxConcurrentRequests)
)

// ProxyHandler is the HTTP handler that proxies requests to the upstream server
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Attempt to acquire a slot in the semaphore, with a timeout for queuing
	select {
	case concurrentRequests <- struct{}{}:
		defer func() { <-concurrentRequests }() // Release the slot when done
		// Forward the request to the upstream server
		proxyRequest(w, r)
	case <-time.After(queueTimeout):
		// If we can't acquire a slot within the timeout, return 503 Service Unavailable
		http.Error(w, "503 Service Unavailable", http.StatusServiceUnavailable)
	}
}

// proxyRequest forwards the request to the upstream server
func proxyRequest(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP request for the upstream server
	req, err := http.NewRequest(r.Method, upstreamURL+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy relevant headers from the incoming request
	req.Header = r.Header

	// Make the request to the upstream server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "502 Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy the upstream response back to the client
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	// Set up a simple HTTP server
	http.HandleFunc("/", ProxyHandler)

	// Start the server
	serverAddr := ":8080"
	fmt.Printf("Starting proxy server on %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
