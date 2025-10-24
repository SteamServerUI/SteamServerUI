package gallery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

// SettingValue represents a setting with a single value type
type SettingValue struct {
	Name        string  `json:"name"`
	IntValue    *int    `json:"intValue,omitempty"`
	BoolValue   *bool   `json:"boolValue,omitempty"`
	StringValue *string `json:"stringValue,omitempty"`
	SupportedOS string  `json:"supported_os,omitempty"`
}

// Plugin represents a plugin with a name and optional OS support
type Plugin struct {
	Name        string `json:"name"`
	SupportedOS string `json:"supported_os,omitempty"`
}

// GalleryRunfile represents a runfile in the gallery
type GalleryRunfile struct {
	Name                string         `json:"name"`
	Filename            string         `json:"filename"`
	Version             string         `json:"version"`
	BackgroundURL       string         `json:"background_url"`
	LogoURL             string         `json:"logo_url"`
	SupportedOS         string         `json:"supported_os"`
	MinVersion          string         `json:"min_version"`
	RecommendedSettings []SettingValue `json:"recommended_settings,omitempty"`
	RecommendedPlugins  []Plugin       `json:"recommended_plugins,omitempty"`
}

// galleryCache stores the parsed and filtered runfile list
var (
	galleryCache []GalleryRunfile
	cacheMutex   sync.Mutex
)

// GetRunfileGallery fetches the list of available runfiles from GitHub Pages
func GetRunfileGallery(forceUpdate bool) ([]GalleryRunfile, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Return cached results if not forcing an update and cache is populated
	if !forceUpdate && len(galleryCache) > 0 {
		logger.Runfile.Debug("Serving runfile gallery from cache")
		return galleryCache, nil
	}

	// Fetch manifest from GitHub Pages
	const manifestURL = "https://steamserverui.github.io/runfiles/manifest.ssui"
	logger.Runfile.Debug("Fetching runfile gallery from " + manifestURL)
	resp, err := http.Get(manifestURL)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to fetch manifest: %v", err))
		return nil, fmt.Errorf("couldn't reach the gallery, network's playing hide and seek")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Runfile.Error(fmt.Sprintf("Manifest fetch failed with status: %d", resp.StatusCode))
		return nil, fmt.Errorf("gallery failed to fetch manifest from GitHub, status: %d", resp.StatusCode)
	}

	// Parse manifest
	var runfiles []GalleryRunfile
	if err := json.NewDecoder(resp.Body).Decode(&runfiles); err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to parse manifest: %v", err))
		return nil, fmt.Errorf("manifest is gibberish, can't make sense of it")
	}

	// Filter by backend version and OS
	currentVersion := config.GetVersion()
	currentOS := runtime.GOOS // "linux", "windows", etc.
	var filtered []GalleryRunfile
	for _, rf := range runfiles {
		if compareVersions(rf.MinVersion, currentVersion) <= 0 {
			// Filter recommended settings and plugins by OS
			var filteredSettings []SettingValue
			for _, setting := range rf.RecommendedSettings {
				if isOSCompatible(setting.SupportedOS, currentOS) {
					filteredSettings = append(filteredSettings, setting)
				}
			}

			var filteredPlugins []Plugin
			for _, plugin := range rf.RecommendedPlugins {
				if isOSCompatible(plugin.SupportedOS, currentOS) {
					filteredPlugins = append(filteredPlugins, plugin)
				}
			}

			// Create a new runfile manifest entry with filtered settings and plugins
			filteredRunfileManifest := rf
			filteredRunfileManifest.RecommendedSettings = filteredSettings
			filteredRunfileManifest.RecommendedPlugins = filteredPlugins
			filtered = append(filtered, filteredRunfileManifest)
		} else {
			logger.Runfile.Debug(fmt.Sprintf("Skipping runfile %s, requires version %s, current is %s", rf.Name, rf.MinVersion, currentVersion))
		}
	}

	// Update cache
	galleryCache = filtered
	logger.Runfile.Info(fmt.Sprintf("Fetched and cached %d runfiles", len(filtered)))

	if len(filtered) == 0 {
		logger.Runfile.Warn("No runfiles compatible with backend version " + currentVersion)
	}

	return filtered, nil
}

// isOSCompatible checks if the supported_os field matches the current OS
func isOSCompatible(supportedOS, currentOS string) bool {
	supportedOS = strings.ToLower(strings.TrimSpace(supportedOS))
	if supportedOS == "" || supportedOS == "all" {
		return true
	}
	return supportedOS == currentOS
}

// saveRunfileToDisk downloads a runfile by identifier and saves it to RunfilesDir
func SaveRunfileToDisk(identifier string, redownload bool) error {
	// Validate identifier: reject if contains path separators or ".."
	if strings.Contains(identifier, "/") || strings.Contains(identifier, "\\") || strings.Contains(identifier, "..") {
		return fmt.Errorf("invalid identifier: path traversal or separator detected")
	}
	filename := fmt.Sprintf("run%s.ssui", identifier)
	baseURL := "https://steamserverui.github.io/runfiles"
	fileURL := fmt.Sprintf("%s/%s", baseURL, filename)

	logger.Runfile.Debug("Fetching runfile from " + fileURL)
	resp, err := http.Get(fileURL)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to fetch runfile %s: %v", filename, err))
		return fmt.Errorf("couldn't grab %s, network's being a jerk", filename)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Runfile.Error(fmt.Sprintf("Runfile %s fetch failed with status: %d", filename, resp.StatusCode))
		return fmt.Errorf("%s is playing hard to get, status: %d", filename, resp.StatusCode)
	}

	saveFilePath := filepath.Join(config.GetRunfilesFolder(), filename)
	logger.Runfile.Debug("Saving runfile to " + saveFilePath)

	// get the dir of the saveFilePath, and os.MkdirAll it if it doesn't exist
	dir := filepath.Dir(saveFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Runfile.Error(fmt.Sprintf("Failed to create runfiles directory %s: %v", dir, err))
			return fmt.Errorf("couldn't create directory")
		}
	}

	// Check if file already exists
	if _, err := os.Stat(saveFilePath); err == nil && !redownload {
		logger.Plugin.Info(fmt.Sprintf("Runfile %s already exists at %s", identifier, saveFilePath))
		return fmt.Errorf("runfile %s already exists as %s", identifier, saveFilePath)
	}

	// Check if file already exists
	if _, err := os.Stat(saveFilePath); err == nil && redownload {
		// create old folder if it doesn't exist
		oldDir := filepath.Join(filepath.Dir(saveFilePath), "old")
		if err := os.MkdirAll(oldDir, 0755); err != nil {
			logger.Runfile.Error(fmt.Sprintf("Failed to create old directory %s: %v", oldDir, err))
			return fmt.Errorf("couldn't create old directory %s, disk's being dramatic", oldDir)
		}

		// create backup filename in old folder
		backupFilename := fmt.Sprintf("%s-old-%s.bak", filepath.Base(saveFilePath), time.Now().Format("2006-01-02_15-04-05"))
		newfp := filepath.Join(oldDir, backupFilename)

		// rename old file to old folder
		if err := os.Rename(saveFilePath, newfp); err != nil {
			logger.Runfile.Error(fmt.Sprintf("Failed to rename old runfile %s to %s: %v", saveFilePath, newfp, err))
			return fmt.Errorf("couldn't rename %s to %s, disk's being dramatic", saveFilePath, newfp)
		}
	}

	// Create the file
	file, err := os.Create(saveFilePath)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to create file %s: %v", saveFilePath, err))
		return fmt.Errorf("disk's throwing a fit, can't save file")
	}
	defer file.Close()

	// Copy response body to file
	if _, err := io.Copy(file, resp.Body); err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to save runfile %s: %v", filename, err))
		return fmt.Errorf("couldn't save %s, disk's being dramatic", filename)
	}

	logger.Runfile.Debug("Successfully saved runfile " + filename)
	loader.InitRunfile(identifier)
	return nil
}

// compareVersions compares two semantic version strings (x.y.z)
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	// Ensure both versions have 3 parts
	for i := 0; i < 3; i++ {
		var n1, n2 int
		if i < len(v1Parts) {
			fmt.Sscanf(v1Parts[i], "%d", &n1)
		}
		if i < len(v2Parts) {
			fmt.Sscanf(v2Parts[i], "%d", &n2)
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	return 0
}
