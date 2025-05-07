package setup

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/steammgr"
)

// BepInEx version: 5.4.23.2 or v5-lts
// SSCM version: 1.0.0

var installMutex sync.Mutex

func CheckAndInstallBepInEx() {
	// Ensure thread safety
	installMutex.Lock()
	defer installMutex.Unlock()

	logger.Install.Info("Checking for BepInEx installation...")

	// Check if BepInEx is already installed
	if _, err := os.Stat("BepInEx"); err == nil {
		logger.Install.Info("BepInEx folder already exists, skipping installation")
		return
	}

	// Determine the URL based on platform
	var url string
	if runtime.GOOS == "windows" {
		url = "https://github.com/BepInEx/BepInEx/releases/download/v5.4.23.2/BepInEx_win_x64_5.4.23.2.zip"
		logger.Install.Info("Detected Windows platform, using Windows BepInEx package")
	} else {
		url = "https://github.com/BepInEx/BepInEx/releases/download/v5.4.23.2/BepInEx_linux_x64_5.4.23.2.zip"
		logger.Install.Info("Detected non-Windows platform, using Linux BepInEx package")
	}

	// Download and install BepInEx
	if err := downloadAndInstallBepInEx(url); err != nil {
		logger.Install.Error(fmt.Sprintf("❌Failed to install BepInEx: %v", err))
	} else {
		logger.Install.Info("✅BepInEx installed successfully")
	}
}

// downloadAndInstallBepInEx downloads the BepInEx zip and extracts it to the current directory
func downloadAndInstallBepInEx(url string) error {
	// Create a temporary file to store the downloaded zip
	tempFile, err := os.CreateTemp("", "bepinex_*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the zip file when done

	// Download the BepInEx zip file
	logger.Install.Info("📥Downloading BepInEx from: " + url)
	err = downloadFile(tempFile.Name(), url)
	if err != nil {
		return fmt.Errorf("failed to download BepInEx: %w", err)
	}

	// Get file info for the zip
	fileInfo, err := tempFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Reopen the file for reading
	tempFile.Close()
	zipFile, err := os.Open(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open downloaded zip: %w", err)
	}
	defer zipFile.Close()

	// Extract the zip file to the current directory
	logger.Install.Info("📦Extracting BepInEx to current directory")
	err = steammgr.Unzip(zipFile, fileInfo.Size(), ".")
	if err != nil {
		return fmt.Errorf("failed to extract BepInEx: %w", err)
	}

	// Clean up changelog.txt if it exists
	if _, err := os.Stat("changelog.txt"); err == nil {
		logger.Install.Debug("🗑️Removing changelog.txt")
		if err := os.Remove("changelog.txt"); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️Failed to remove changelog.txt: %v", err))
		}
	}

	if runtime.GOOS == "linux" {
		// make sure run_bepinex.sh is executable
		if err := os.Chmod("./run_bepinex.sh", os.ModePerm); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️Failed to make run_bepinex.sh executable: %v", err))
		}
	}

	return nil
}

func InstallSSCM() {
	logger.Install.Info("🕑Installing SSCM...")

	CheckAndInstallBepInEx()

	// Enable SSCM
	config.SetIsSSCMEnabled(true)

	logger.Install.Info("✅SSCM enabled")
}
