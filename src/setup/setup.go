package setup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/steammgr"
)

var (
	IsSetupComplete bool
	V6setupMutex    sync.Mutex
)

var downloadBranch string // Holds the branch to download from

// Install performs the entire installation process and ensures the server waits for it to complete
func Install(wg *sync.WaitGroup) {
	defer wg.Done() // Signal that installation is complete

	// Step 0: Check for updates
	if err := UpdateExecutable(); err != nil {
		logger.Install.Error("❌Update check went sideways: " + err.Error())
	}

	// Step 1: Check and download the UIMod folder contents
	logger.Install.Info("🔄Checking UIMod folder contents...")
	CheckAndDownloadUIMod()
	logger.Install.Info("✅UIMod folder setup complete.")
	time.Sleep(2 * time.Second) // Small pause to let the user read potential errors
	// Step 3: Install and run SteamCMD
	logger.Install.Info("🔄Installing SteamCMD...")
	steammgr.InstallSteamCMD()
	logger.Install.Info("✅Setup complete!")
	V6setupMutex.Lock()
	IsSetupComplete = true
	V6setupMutex.Unlock()
}

// fileMappings defines the mapping of local file paths to their GitHub raw URLs with a {branch} placeholder
var fileMappings = map[string]string{
	// v2 UI
	"v2/index.html":      "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/v2/index.html",
	"v2/assets/ssui.css": "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/v2/assets/ssui.css",
	"v2/assets/ssui.js":  "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/v2/assets/ssui.js",
	// v1 UI
	"ui/config.html":                  "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/ui/config.html",
	"ui/index.html":                   "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/ui/index.html",
	"ui/detectionmanager.html":        "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/ui/detectionmanager.html",
	"assets/stationeers.png":          "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/stationeers.png",
	"assets/favicon.ico":              "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/favicon.ico",
	"assets/apiinfo.html":             "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/apiinfo.html",
	"twoboxform/twoboxform.css":       "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/twoboxform/twoboxform.css",
	"twoboxform/twoboxform.js":        "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/twoboxform/twoboxform.js",
	"twoboxform/twoboxform.html":      "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/twoboxform/twoboxform.html",
	"assets/css/apiinfo.css":          "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/apiinfo.css",
	"assets/css/background.css":       "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/background.css",
	"assets/css/base.css":             "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/base.css",
	"assets/css/components.css":       "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/components.css",
	"assets/css/config.css":           "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/config.css",
	"assets/css/detectionmanager.css": "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/detectionmanager.css",
	"assets/css/home.css":             "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/home.css",
	"assets/css/mobile.css":           "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/mobile.css",
	"assets/css/runfileterminal.css":  "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/runfileterminal.css",
	"assets/css/style.css":            "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/style.css",
	"assets/css/tabs.css":             "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/tabs.css",
	"assets/css/variables.css":        "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/css/variables.css",
	"assets/js/main.js":               "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/main.js",
	"assets/js/detectionmanager.js":   "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/detectionmanager.js",
	"assets/js/console-manager.js":    "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/console-manager.js",
	"assets/js/server-api.js":         "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/server-api.js",
	"assets/js/ui-utils.js":           "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/ui-utils.js",
	"assets/js/runfile-terminal.js":   "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/runfile-terminal.js",
	"assets/js/runfile-settings.js":   "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/runfile-settings.js",
	"assets/js/settings.js":           "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/assets/js/settings.js",
}

// CheckAndDownloadUIMod ensures the UI module is present and up-to-date
func CheckAndDownloadUIMod() {
	// Define directories
	uiModDir := config.GetUIModFolder()
	dirs := []string{
		uiModDir,
		uiModDir + "twoboxform/",
		uiModDir + "detectionmanager/",
		uiModDir + "assets/",
		uiModDir + "assets/css/",
		uiModDir + "ui/",
		uiModDir + "config/",
		uiModDir + "tls/",
		uiModDir + "assets/js/",
		config.GetRunFilesFolder(),
		uiModDir + "v2/",
		uiModDir + "v2/assets/",
	}

	// Determine the branch to download from
	downloadBranch := config.Branch
	if config.Branch == "release" || config.Branch == "Release" {
		downloadBranch = "main"
	}
	logger.Install.Info("Using branch: " + downloadBranch)

	// Create resolved file mappings with the correct branch
	resolvedFiles := make(map[string]string)
	for relativePath, urlTemplate := range fileMappings {
		localPath := uiModDir + relativePath
		remoteURL := strings.Replace(urlTemplate, "{branch}", downloadBranch, 1)
		resolvedFiles[localPath] = remoteURL
	}

	// check if UIMod exists
	uiModExists := true
	if _, err := os.Stat(uiModDir); os.IsNotExist(err) {
		uiModExists = false
		logger.Install.Warn("🍲 Unable to find UIMod folder. Cooking it...")
	}

	// Always ensure all directories exist
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				logger.Install.Error("❌ Error occurred while creating the folder structure: " + err.Error())
				return
			}
			logger.Install.Debug("⚠️ Created folder: " + dir)
		}
	}

	// Then decide whether to download all files or just update
	if !uiModExists {
		downloadAllFiles(resolvedFiles)
	} else {
		// Directory exists
		logger.Install.Debug(fmt.Sprintf("IsUpdateEnabled: %v", config.GetIsUpdateEnabled()))
		if config.GetIsUpdateEnabled() {
			logger.Install.Info("🔍 Validating UIMod files for updates...")
			updateFilesIfDifferent(resolvedFiles)
		} else {
			logger.Install.Info("♻️ Folder ./UIMod already exists. Updates disabled, skipping validation.")
		}
	}
}

// downloadFile downloads a file from the given URL to the specified downloadFilePath
func downloadFile(downloadFilePath, url string) error {
	// Ensure the directory exists
	dir := filepath.Dir(downloadFilePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	// Create the file
	out, err := os.Create(downloadFilePath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Fetch the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// Write to file
	_, err = io.Copy(out, resp.Body)
	return err
}
