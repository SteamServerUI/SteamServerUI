package pluginsapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

// a http handler that can be used to use logger.Plugin from the api
func PluginLogHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Only handle POST requests
	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","message":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Define struct for incoming JSON
	type logRequest struct {
		Level      string `json:"level"`
		PluginName string `json:"pluginname"`
		Message    string `json:"message"`
	}

	// Decode JSON request body
	var req logRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, `{"status":"error","message":"Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Level == "" || req.PluginName == "" || req.Message == "" {
		http.Error(w, `{"status":"error","message":"Missing required fields"}`, http.StatusBadRequest)
		return
	}

	msg := "<" + req.PluginName + "> " + req.Message
	// Log based on level
	switch strings.ToLower(req.Level) {
	case "info":
		logger.Plugin.Info(msg)
	case "warn":
		logger.Plugin.Warn(msg)
	case "error":
		logger.Plugin.Error(msg)
	case "debug":
		logger.Plugin.Debug(msg)
	default:
		http.Error(w, `{"status":"error","message":"Invalid log level"}`, http.StatusBadRequest)
		return
	}

	// Write success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
