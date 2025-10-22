package pluginproxy

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

func UnixSocketProxyHandler(socketPath string, pluginName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Trim the plugin prefix from the URL path
		subPath := strings.TrimPrefix(r.URL.Path, "/plugins/"+pluginName)
		if subPath == "" {
			subPath = "/" // Default to root if no subpath
		}
		requestPath := subPath
		if r.URL.RawQuery != "" {
			requestPath += "?" + r.URL.RawQuery
		}

		// Dial the Unix domain socket
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			logger.Plugin.Debugf("Failed to connect to Plugin socket: %v", err)
			http.Error(w, "Failed to connect to Plugin socket, the plugin likely crashed, was stopped or is unhealty. To remove plugin fully, restart the Backend. "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		var requestBuilder strings.Builder
		requestBuilder.WriteString(r.Method + " " + requestPath + " HTTP/1.1\r\n")
		requestBuilder.WriteString("Host: localhost\r\n")

		// Copy headers from the original request
		for key, values := range r.Header {
			for _, value := range values {
				requestBuilder.WriteString(key + ": " + value + "\r\n")
			}
		}
		requestBuilder.WriteString("\r\n")

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Plugin.Debugf("Failed to read request body: %v", err)
			http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		requestBuilder.Write(body)

		_, err = conn.Write([]byte(requestBuilder.String()))
		if err != nil {
			logger.Plugin.Debugf("Failed to write to Unix socket: %v", err)
			http.Error(w, "Failed to write to Unix socket: "+err.Error(), http.StatusInternalServerError)
			return
		}

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
