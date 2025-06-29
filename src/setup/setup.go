package setup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	// NOTE: Currently empty as files are now embedded in the executable. Keep this structure for future use.

	// v1 UI - commented out since files are embedded, left here for reference in case we need this funcitonality again
	// "ui/config.html":           "https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/{branch}/UIMod/ui/config.html",
}

// CheckAndDownloadUIMod ensures the UI module is present and up-to-date
func CheckAndDownloadUIMod() {
	// Define directories
	uiModDir := config.GetUIModFolder()
	dirs := []string{
		uiModDir,
		uiModDir + "config/",
		uiModDir + "config/tls/",
		config.GetRunFilesFolder(),
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

	// Check if fileMappings is empty - skip download if so
	if len(fileMappings) == 0 {
		logger.Install.Debug("📁 File mappings empty - no additional files to download available")
		return
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
