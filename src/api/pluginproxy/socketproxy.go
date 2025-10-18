package pluginproxy

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

func UnixSocketProxyHandler(socketPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract the subpath after /plugins/ExamplePlugin/
		subPath := strings.TrimPrefix(r.URL.Path, "/plugins/ExamplePlugin")
		if subPath == "" {
			subPath = "/" // Default to root if no subpath
		}

		// Dial the Unix domain socket
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			logger.Plugin.Debugf("Failed to connect to Unix socket: %v", err)
			http.Error(w, "Failed to connect to Unix socket: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Construct a full HTTP request with the subpath
		var requestBuilder strings.Builder
		requestBuilder.WriteString(r.Method + " " + subPath + " HTTP/1.1\r\n")
		requestBuilder.WriteString("Host: localhost\r\n")
		requestBuilder.WriteString("\r\n")

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Plugin.Debugf("Failed to read request body: %v", err)
			http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		requestBuilder.Write(body)

		// Send the HTTP request to the Unix socket
		_, err = conn.Write([]byte(requestBuilder.String()))
		if err != nil {
			logger.Plugin.Debugf("Failed to write to Unix socket: %v", err)
			http.Error(w, "Failed to write to Unix socket: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Read the response from the Unix socket
		reader := bufio.NewReader(conn)
		resp, err := http.ReadResponse(reader, r)
		if err != nil {
			logger.Plugin.Debugf("Failed to read response from Unix socket: %v", err)
			http.Error(w, "Failed to read response from Unix socket: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy headers from the socket response
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Set the status code
		w.WriteHeader(resp.StatusCode)

		// Copy the response body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			logger.Plugin.Debugf("Failed to write response to client: %v", err)
		}
	}
}
