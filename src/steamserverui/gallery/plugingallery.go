package gallery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

// GalleryPlugin represents a plugin in the gallery
type GalleryPlugin struct {
	Name          string `json:"name"`
	Filename      string `json:"filename"`
	Version       string `json:"version"`
	BackgroundURL string `json:"background_url"`
	LogoURL       string `json:"logo_url"`
	SupportedOS   string `json:"supported_os"`
	MinVersion    string `json:"min_version"`
}

// pluginCache stores the parsed and filtered plugin list
var (
	pluginCache      []GalleryPlugin
	pluginCacheMutex sync.Mutex
)

// GetPluginGallery fetches the list of available plugins from GitHub Pages
func GetPluginGallery(forceUpdate bool) ([]GalleryPlugin, error) {
	pluginCacheMutex.Lock()
	defer pluginCacheMutex.Unlock()

	// Return cached results if not forcing an update and cache is populated
	if !forceUpdate && len(pluginCache) > 0 {
		logger.Plugin.Debug("Serving plugin gallery from cache")
		return pluginCache, nil
	}

	// Fetch manifest from GitHub Pages
	const manifestURL = "https://steamserverui.github.io/plugins/manifest.ssui"
	logger.Plugin.Debug("Fetching plugin gallery from " + manifestURL)
	resp, err := http.Get(manifestURL)
	if err != nil {
		logger.Plugin.Error(fmt.Sprintf("Failed to fetch plugin manifest: %v", err))
		return nil, fmt.Errorf("couldn't reach the plugin gallery, network's playing hide and seek")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Plugin.Error(fmt.Sprintf("Plugin manifest fetch failed with status: %d", resp.StatusCode))
		return nil, fmt.Errorf("plugin gallery failed to fetch manifest from GitHub, status: %d", resp.StatusCode)
	}

	// Parse manifest
	var plugins []GalleryPlugin
	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		logger.Plugin.Error(fmt.Sprintf("Failed to parse plugin manifest: %v", err))
		return nil, fmt.Errorf("plugin manifest is gibberish, can't make sense of it")
	}

	// Filter by backend version and OS
	currentVersion := config.GetVersion()
	currentOS := runtime.GOOS
	var filtered []GalleryPlugin
	for _, p := range plugins {
		// Check version compatibility
		if compareVersions(p.MinVersion, currentVersion) > 0 {
			logger.Plugin.Debug(fmt.Sprintf("Skipping plugin %s, requires version %s, current is %s", p.Name, p.MinVersion, currentVersion))
			continue
		}
		// Check OS compatibility
		if p.SupportedOS != "all" && p.SupportedOS != currentOS {
			logger.Plugin.Debug(fmt.Sprintf("Skipping plugin %s, supported OS %s, current OS is %s", p.Name, p.SupportedOS, currentOS))
			continue
		}
		filtered = append(filtered, p)
	}

	// Update cache
	pluginCache = filtered
	logger.Plugin.Info(fmt.Sprintf("Fetched and cached %d plugins", len(filtered)))

	if len(filtered) == 0 {
		logger.Plugin.Warn("No plugins compatible with backend version " + currentVersion + " and OS " + currentOS)
	}

	return filtered, nil
}

// SavePluginToDisk downloads a plugin by name and saves it to PluginsDir
func SavePluginToDisk(name string, redownload bool) error {
	// Find plugin in cache to get filename
	pluginCacheMutex.Lock()
	var plugin GalleryPlugin
	found := false
	for _, p := range pluginCache {
		if p.Name == name {
			plugin = p
			found = true
			break
		}
	}
	pluginCacheMutex.Unlock()

	if !found {
		logger.Plugin.Error(fmt.Sprintf("Plugin %s not found in gallery", name))
		return fmt.Errorf("plugin %s not found in gallery", name)
	}

	filename := plugin.Filename
	baseURL := "https://steamserverui.github.io/plugins"
	fileURL := fmt.Sprintf("%s/%s", baseURL, filename)

	saveFilePath := filepath.Join(config.GetPluginsFolder(), filename)

	// Check if file already exists
	if _, err := os.Stat(saveFilePath); err == nil && !redownload {
		logger.Plugin.Info(fmt.Sprintf("Plugin %s already exists at %s", name, saveFilePath))
		return fmt.Errorf("plugin %s already exists", name)
	}

	logger.Plugin.Debug("Fetching plugin from " + fileURL)
	resp, err := http.Get(fileURL)
	if err != nil {
		logger.Plugin.Error(fmt.Sprintf("Failed to fetch plugin %s: %v", filename, err))
		return fmt.Errorf("couldn't grab %s, network's being a jerk", filename)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Plugin.Error(fmt.Sprintf("Plugin %s fetch failed with status: %d", filename, resp.StatusCode))
		return fmt.Errorf("%s is playing hard to get, status: %d", filename, resp.StatusCode)
	}

	logger.Plugin.Debug("Saving plugin to " + saveFilePath)

	// Create directory if it doesn't exist
	dir := filepath.Dir(saveFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Plugin.Error(fmt.Sprintf("Failed to create plugins directory %s: %v", dir, err))
			return fmt.Errorf("couldn't create directory")
		}
	}

	// Create file if it doesn't exist
	file, err := os.Create(saveFilePath)
	if err != nil {
		logger.Plugin.Error(fmt.Sprintf("Failed to create file %s: %v", saveFilePath, err))
		return fmt.Errorf("disk's throwing a fit, can't save file")
	}

	// Copy response body to file
	if _, err := io.Copy(file, resp.Body); err != nil {
		file.Close()
		logger.Plugin.Error(fmt.Sprintf("Failed to save plugin %s: %v", filename, err))
		return fmt.Errorf("couldn't save %s, disk's being dramatic", filename)
	}

	// Flush writes to disk and close file
	if err := file.Sync(); err != nil {
		file.Close()
		logger.Plugin.Error(fmt.Sprintf("Failed to sync plugin file %s: %v", saveFilePath, err))
		return fmt.Errorf("couldn't sync %s to disk", filename)
	}
	if err := file.Close(); err != nil {
		logger.Plugin.Error(fmt.Sprintf("Failed to close plugin file %s: %v", saveFilePath, err))
		return fmt.Errorf("couldn't close %s", filename)
	}

	// Brief delay to ensure filesystem releases the file
	time.Sleep(100 * time.Millisecond)

	// Set executable permissions on Linux
	if runtime.GOOS == "linux" {
		if err := os.Chmod(saveFilePath, 0755); err != nil {
			logger.Plugin.Error(fmt.Sprintf("Failed to set executable permissions for %s: %v", saveFilePath, err))
			return fmt.Errorf("couldn't make %s executable, permissions issue", filename)
		}
	}

	logger.Plugin.Debug("Successfully saved plugin " + filename)
	// add to config.SetRegisteredPlugins map string string by pluginname, filename
	config.SetRegisteredPlugins(map[string]string{name: filename})
	loader.ReloadBackend()
	return nil
}
