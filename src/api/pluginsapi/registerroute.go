package pluginsapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pluginproxy"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

func RegisterPluginRouteHandler(w http.ResponseWriter, r *http.Request, protectedMux *http.ServeMux) {
	w.Header().Set("Content-Type", "application/json")

	// Only handle POST requests
	if r.Method != http.MethodPost {
		http.Error(w, `{"status":"error","message":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Define struct for incoming JSON
	type registerRequest struct {
		PluginName string `json:"pluginname"`
	}

	// Decode JSON request body
	var req registerRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, `{"status":"error","message":"Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.PluginName == "" {
		http.Error(w, `{"status":"error","message":"Missing required field pluginname"}`, http.StatusBadRequest)
		return
	}

	logger.API.Infof("Registering plugin in API: %s", req.PluginName)

	// Dynamically register the plugin route in protectedMux
	route := fmt.Sprintf("/plugins/%s/", req.PluginName)
	socketPath := fmt.Sprintf("/tmp/ssui/%s.sock", req.PluginName)
	protectedMux.HandleFunc(route, pluginproxy.UnixSocketProxyHandler(socketPath))

	// Write success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
