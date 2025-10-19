package pluginsapi

import (
	"encoding/json"
	"net/http"
)

func HandleListPlugins(w http.ResponseWriter, r *http.Request) {
	plugins := make([]string, 0)
	for plugin := range pluginRoutes {
		plugins = append(plugins, plugin)
	}
	json.NewEncoder(w).Encode(plugins)
}
