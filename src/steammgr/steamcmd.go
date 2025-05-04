package steammgr

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/SteamServerUI/SteamServerUI/v6/src/argmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

// ExtractorFunc is a type that represents a function for extracting archives.
// It takes an io.ReaderAt, the size of the content, and the destination directory.
type ExtractorFunc func(io.ReaderAt, int64, string) error

// Constants for repeated strings
var (
	SteamCMDLinuxURL   = config.GetSteamCMDLinuxURL()
	SteamCMDWindowsURL = config.GetSteamCMDWindowsURL()
	SteamCMDLinuxDir   = config.GetSteamCMDLinuxDir()
	SteamCMDWindowsDir = config.GetSteamCMDWindowsDir()
	RuntimeSteamCMDDir string
)

// InstallAndRunSteamCMD installs and runs SteamCMD based on the platform (Windows/Linux).
// It automatically detects the OS and calls the appropriate installation function.
func InstallSteamCMD() {
	if runtime.GOOS == "windows" {
		installSteamCMDWindows()
	} else if runtime.GOOS == "linux" {
		installSteamCMDLinux()
	} else {
		logger.Install.Error("❌ SteamCMD installation is not supported on this OS.\n")
		return
	}
}

func installSteamCMD(platform string, steamCMDDir string, downloadURL string, extractFunc ExtractorFunc) {
	// Check if SteamCMD is already installed
	if _, err := os.Stat(steamCMDDir); os.IsNotExist(err) {
		logger.Install.Warn("⚠️ SteamCMD not found for " + platform + ", downloading...\n")

		// Create SteamCMD directory
		if err := createSteamCMDDirectory(steamCMDDir); err != nil {
			logger.Install.Error("❌ Error creating SteamCMD directory: " + err.Error() + "\n")
			return
		}

		// Ensure cleanup on failure
		success := false
		defer func() {
			if !success {
				logger.Install.Warn("⚠️ Cleaning up due to failure...\n")
				os.RemoveAll(steamCMDDir)
			}
		}()

		// Install required libraries
		if err := installRequiredLibraries(); err != nil {
			logger.Install.Error("❌ Error installing required libraries: " + err.Error() + "\n")
			return
		}

		// Download and extract SteamCMD
		if err := downloadAndExtractSteamCMD(downloadURL, steamCMDDir, extractFunc); err != nil {
			logger.Install.Error("❌ " + err.Error() + "\n")
			return
		}

		// Set executable permissions for SteamCMD files
		if err := setExecutablePermissions(steamCMDDir); err != nil {
			logger.Install.Error("❌ Error setting executable permissions: " + err.Error() + "\n")
			return
		}

		// Verify the steamcmd binary
		if err := verifySteamCMDBinary(steamCMDDir); err != nil {
			logger.Install.Error("❌ " + err.Error() + "\n")
			return
		}

		// Mark installation as successful
		success = true
		logger.Install.Info("✅ SteamCMD installed successfully.\n")
	} else {

		logger.Install.Info("✅ SteamCMD is already installed.\n")
	}

	//RunSteamCMD() // disabled in v6, steamcmd is now triggered after a game is defined.
}

// installSteamCMDLinux downloads and installs SteamCMD on Linux.
func installSteamCMDLinux() {
	RuntimeSteamCMDDir = SteamCMDLinuxDir
	installSteamCMD("Linux", SteamCMDLinuxDir, SteamCMDLinuxURL, untarWrapper)
}

// installSteamCMDWindows downloads and installs SteamCMD on Windows.
func installSteamCMDWindows() {
	RuntimeSteamCMDDir = SteamCMDWindowsDir
	installSteamCMD("Windows", SteamCMDWindowsDir, SteamCMDWindowsURL, Unzip)
}

// runSteamCMD runs the SteamCMD command to update the game.
func RunSteamCMD() {
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Install.Error("❌ Error getting current working directory: " + err.Error() + "\n")
		return
	}
	logger.Install.Debug("✅ Current working directory: " + currentDir + "\n")

	// Ensure permissions every time if we run on linux
	if runtime.GOOS != "windows" {
		if err := setExecutablePermissions(RuntimeSteamCMDDir); err != nil {
			logger.Install.Error("❌ Error setting steamcmd executable permissions, your Steamcmd install might be broken: " + err.Error() + "\n")
			return
		}
	}

	// Build SteamCMD command
	cmd := buildSteamCMDCommand(currentDir)

	// Set output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	logger.Install.Info("🕑 Running SteamCMD...\n")
	err = cmd.Run()
	if err != nil {
		logger.Install.Error("❌ Error running SteamCMD: " + err.Error() + "\n")
		return
	}
	logger.Install.Info("✅ SteamCMD executed successfully.\n")
}

// buildSteamCMDCommand constructs the SteamCMD command based on the OS.
func buildSteamCMDCommand(currentDir string) *exec.Cmd {
	//print the config.GetGameBranch() and config.GetGameServerAppID()
	logger.Install.Info("🔍 SSUI Server Identifier: " + argmgr.CurrentRunfile.Meta.Name + "\n")
	logger.Install.Info("🔍 Game Branch: " + config.GetGameBranch() + "\n")
	logger.Install.Info("🔍 Game Server App ID: " + argmgr.CurrentRunfile.SteamAppID + "\n")
	runfileAppID := argmgr.CurrentRunfile.SteamAppID

	if runtime.GOOS == "windows" {
		return exec.Command(filepath.Join(RuntimeSteamCMDDir, "steamcmd.exe"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", runfileAppID, "-beta", config.GetGameBranch(), "validate", "+quit")
	}

	if config.GetGameBranch() == "public" {
		return exec.Command(filepath.Join(RuntimeSteamCMDDir, "steamcmd.sh"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", runfileAppID, "validate", "+quit")
	}
	return exec.Command(filepath.Join(RuntimeSteamCMDDir, "steamcmd.sh"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", runfileAppID, "-beta", config.GetGameBranch(), "validate", "+quit")
}
