package web

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

// HandleCodeServer proxies requests from /api/v2/codeserver to the code-server Unix socket.
func HandleCodeServer(w http.ResponseWriter, r *http.Request) {

	if !config.GetIsCodeServerEnabled() {
		http.Error(w, "Code-server is not enabled", http.StatusServiceUnavailable)
		return
	}
	// Define the correct Unix socket path (relative to current directory).
	var codeServerSocket = config.CodeServerSocketPath

	// Check if the socket exists to avoid unnecessary dial attempts.
	if _, err := os.Stat(codeServerSocket); os.IsNotExist(err) {
		http.Error(w, "Code-server is not running", http.StatusServiceUnavailable)
		fmt.Printf("Proxy error: socket %s not found: %v\n", codeServerSocket, err)
		return
	}

	// Dial the Unix socket.
	conn, err := net.Dial("unix", codeServerSocket)
	if err != nil {
		http.Error(w, "Failed to connect to code-server", http.StatusServiceUnavailable)
		fmt.Printf("Proxy error: failed to dial socket %s: %v\n", codeServerSocket, err)
		return
	}
	defer conn.Close()

	// Check for WebSocket upgrade request.
	if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		strings.ToLower(r.Header.Get("Connection")) == "upgrade" {
		proxyWebSocket(w, r, conn)
		return
	}
	urlPath := strings.TrimPrefix(r.URL.Path, "/api/v2/codeserver")
	if urlPath == "" {
		urlPath = "/"
	}

	// Create a proper URL for the Unix socket request
	// The host doesn't matter since we're overriding the transport
	fullURL := "http://interalcodeserver" + urlPath
	if r.URL.RawQuery != "" {
		fullURL += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, fullURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		fmt.Printf("Proxy error: failed to create request: %v\n", err)
		return
	}

	// Copy headers from the original request.
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Use a custom HTTP client to send the request to the Unix socket.
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("unix", codeServerSocket)
			},
		},
	}

	// Send the request to code-server.
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to proxy to code-server", http.StatusServiceUnavailable)
		fmt.Printf("Proxy error: failed to proxy request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Copy response headers and status to the client.
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Copy response body to the client.
	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Printf("Proxy error: failed to copy response: %v\n", err)
	}
}

// proxyWebSocket handles WebSocket connections to code-server's Unix socket.
func proxyWebSocket(w http.ResponseWriter, r *http.Request, conn net.Conn) {
	// Hijack the HTTP connection to manage WebSocket manually.
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Server does not support WebSocket", http.StatusInternalServerError)
		fmt.Println("Proxy error: server does not support hijacking")
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Failed to hijack connection", http.StatusInternalServerError)
		fmt.Printf("Proxy error: failed to hijack connection: %v\n", err)
		return
	}
	defer clientConn.Close()

	// Send WebSocket handshake headers to code-server.
	handshake := fmt.Sprintf(
		"GET %s HTTP/1.1\r\nHost: unix\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n",
		r.URL.String(),
		r.Header.Get("Sec-WebSocket-Key"),
	)
	if _, err := conn.Write([]byte(handshake)); err != nil {
		fmt.Printf("Proxy error: failed to send WebSocket handshake: %v\n", err)
		return
	}

	// Read handshake response from code-server.
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Proxy error: failed to read WebSocket handshake: %v\n", err)
		return
	}
	// Write handshake response to client.
	if _, err := clientConn.Write(buf[:n]); err != nil {
		fmt.Printf("Proxy error: failed to write WebSocket handshake to client: %v\n", err)
		return
	}

	// Bidirectional proxying: client <-> code-server.
	errChan := make(chan error, 2)
	go func() {
		// Copy from client to code-server.
		_, err := io.Copy(conn, clientConn)
		errChan <- err
	}()
	go func() {
		// Copy from code-server to client.
		_, err := io.Copy(clientConn, conn)
		errChan <- err
	}()

	// Wait for an error or connection close.
	<-errChan
}
