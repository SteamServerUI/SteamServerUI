package pluginsapi

import (
	"encoding/json"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/plugins"
)

func HandleListPluginAPIRoutes(w http.ResponseWriter, r *http.Request) {
	plugins := make([]string, 0)
	for plugin := range pluginRoutes {
		plugins = append(plugins, plugin)
	}
	json.NewEncoder(w).Encode(plugins)
}

func HandleListPluginNames(w http.ResponseWriter, r *http.Request) {
	resp := make([]string, 0)

	plugins.RunningPluginsMutex.Lock()
	defer plugins.RunningPluginsMutex.Unlock()
	for plugin := range plugins.RunningPlugins {
		resp = append(resp, plugin)
	}
	json.NewEncoder(w).Encode(resp)
}

func HandleStopPlugin(w http.ResponseWriter, r *http.Request) {

	type stopRequest struct {
		PluginName string `json:"pluginname"`
	}

	var req stopRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, `{"status":"error","message":"Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	pluginname := req.PluginName

	if pluginname == "" {
		http.Error(w, `{"status":"error","message":"Missing required field pluginname"}`, http.StatusBadRequest)
		return
	}

	if err := plugins.StopPlugin(pluginname); err != nil {
		http.Error(w, `{"status":"error","message":"Plugin not running"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Plugin stopped successfully"})
}
