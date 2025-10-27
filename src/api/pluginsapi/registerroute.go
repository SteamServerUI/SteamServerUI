package pluginsapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pluginproxy"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

var pluginRoutes = make(map[string]bool)
var pluginRoutesMu sync.Mutex

func RegisterPluginRouteHandler(w http.ResponseWriter, r *http.Request, apiMux *http.ServeMux, webserverMux *http.ServeMux) {
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
	// sanatize plugin name (allow alphanumeric, underscores, and hyphens only)
	if !isValidPluginName(req.PluginName) {
		http.Error(w, `{"status":"error","message":"Invalid plugin name. Use only alphanumeric characters, underscores, or hyphens"}`, http.StatusBadRequest)
		return
	}
	socketPath := ""

	route := fmt.Sprintf("/plugins/%s/", req.PluginName)

	err := checkRoute(route)
	if err {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Plugin route was already registered"})
		return
	}

	if runtime.GOOS == "linux" {
		socketPath = fmt.Sprintf("./SSUI/plugins/sockets/%s.sock", req.PluginName)

		// check if the plugin socket exists
		if !pluginSocketExists(socketPath) {
			w.WriteHeader(http.StatusNotImplemented)
			json.NewEncoder(w).Encode(map[string]string{"status": "failed", "message": "Plugin socket does not exist. Make sure to call PluginLib.ExposeAPI before calling PluginLib.RegisterPluginAPI"})
			return
		}
		webserverMux.HandleFunc(route, pluginproxy.UnixSocketProxyHandler(socketPath, req.PluginName))
	}

	if runtime.GOOS == "windows" {
		parentpath := getpluginPipeParentPath()
		socketPath = parentpath + req.PluginName

		logger.Socket.Info("HERE I would handle the plugin from " + socketPath + " in the pluginproxy.WindowsNamedPipeProxyHandler")
		//webserverMux.HandleFunc(route, pluginproxy.WindowsNamedPipeProxyHandler(socketPath, req.PluginName))
	}

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

func isValidPluginName(name string) bool {
	// Allow alphanumeric, underscores, and hyphens (minimum 1 character, maximum 50 characters)
	pattern := `^[a-zA-Z0-9_-]{1,50}$`
	matched, err := regexp.MatchString(pattern, name)
	return err == nil && matched
}

func pluginSocketExists(socketPath string) bool {
	_, err := os.Stat(socketPath)
	return err == nil
}

var pluginPipeParentPath string

func getpluginPipeParentPath() string {
	if pluginPipeParentPath != "" {
		return pluginPipeParentPath
	}
	file, err := os.Open("./SSUI/plugins/sockets/pipename.identifier")
	if err != nil {
		fmt.Println("Error opening pipename.identifier file, I have to go...:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pluginPipeParentPath = scanner.Text()
		fmt.Println("Read PluginPipePath:", pluginPipeParentPath) // Debug
	}
	if pluginPipeParentPath == "" {
		fmt.Println("Error reading pipename.identifier file, I have to go...:", err)
		os.Exit(1)
	}
	fmt.Println("Final plugin pipe parent path:", pluginPipeParentPath) // Debug
	return pluginPipeParentPath
}
