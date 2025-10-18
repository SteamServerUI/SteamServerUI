package pluginsapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pluginproxy"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

var pluginRoutes = make(map[string]bool)
var pluginRoutesMu sync.Mutex

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

	// Dynamically register the plugin route in protectedMux
	route := fmt.Sprintf("/plugins/%s/", req.PluginName)
	socketPath := fmt.Sprintf("/tmp/ssui/%s.sock", req.PluginName)

	err := checkRoute(route)
	if err {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"status": "failed", "message": "Plugin route already registered"})
		return
	}

	protectedMux.HandleFunc(route, pluginproxy.UnixSocketProxyHandler(socketPath))
	logger.Plugin.Infof("Registered %s plugin route %s in API", req.PluginName, route)

	// Write success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Plugin route registered successfully"})
}

func checkRoute(route string) (registered bool) {
	// Check if the route is already registered
	pluginRoutesMu.Lock()
	defer pluginRoutesMu.Unlock()

	if pluginRoutes[route] {
		return true
	}
	// save the route in the plugin routes map
	pluginRoutes[route] = true
	return false
}
