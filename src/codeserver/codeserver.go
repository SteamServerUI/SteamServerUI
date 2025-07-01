package codeserver

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

var (
	codeServerInstallDir    = config.CodeServerInstallDir
	codeServerPath          = config.CodeServerPath
	codeServerBinaryPath    = config.CodeServerBinaryPath
	codeServerSocketPath    = config.CodeServerSocketPath
	installScriptURL        = config.CodeServerInstallScriptURL
	configFilePath          = config.CodeServerConfigFilePath
	codeServerUserDataDir   = config.CodeServerUserDataDir
	codeServerExtensionsDir = config.CodeServerExtensionsDir
	settingsFilePath        = config.CodeServerSettingsFilePath
)

// InitCodeServer initializes code-server at server startup.
// Creates the directory, installs, and starts code-server.
func InitCodeServer() error {
	if !config.GetIsCodeServerEnabled() {
		return nil
	}

	os.RemoveAll(codeServerSocketPath)
	// Create directory if it doesn't exist.
	if err := os.MkdirAll(codeServerPath, 0755); err != nil {
		return fmt.Errorf("failed to create Code Server directory: %v", err)
	}

	logger.Main.Info("Initializing Code Server...")
	msg := DownloadInstallCodeServer()
	logger.Main.Info(msg)
	if !strings.Contains(strings.ToLower(msg), "successfully") && !strings.Contains(strings.ToLower(msg), "already installed") {
		return fmt.Errorf("code-server installation failed: %s", msg)
	}

	logger.Main.Info("Starting Code Server...")
	// Start code-server.
	if err := StartCodeServer(); err != nil {
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	return nil
}

// DownloadInstallCodeServer downloads and installs code-server using the official install script.
func DownloadInstallCodeServer() string {
	// Enforce Linux-only support.
	if strings.ToLower(runtime.GOOS) != "linux" {
		return "Code Server is only supported on Linux"
	}

	// Check if code-server binary already exists to avoid re-installing.
	if _, err := os.Lstat(codeServerBinaryPath); err == nil {
		// Verify the symlink target exists.
		target, err := os.Readlink(codeServerBinaryPath)
		if err != nil {
			return fmt.Sprintf("Failed to read symlink %s: %v", codeServerBinaryPath, err)
		}
		if _, err := os.Stat(target); err == nil {
			return "Code Server already installed"
		}
	}

	// Check if code-server binary already exists to avoid re-installing.
	if _, err := os.Stat(codeServerBinaryPath); err == nil {
		return "Code Server already installed"
	}

	// Create a temporary file for the install script.
	tempScript := filepath.Join(codeServerPath, "install.sh")
	if err := os.MkdirAll(filepath.Dir(tempScript), 0755); err != nil {
		return fmt.Sprintf("Failed to create cs directory: %v", err)
	}

	// Download the install script.
	resp, err := http.Get(installScriptURL)
	if err != nil {
		return fmt.Sprintf("Failed to download install script: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Failed to download install script: HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(tempScript)
	if err != nil {
		return fmt.Sprintf("Failed to create temp script: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Sprintf("Failed to save install script: %v", err)
	}

	// Set executable permissions for the script.
	if err := os.Chmod(tempScript, 0755); err != nil {
		return fmt.Sprintf("Failed to set script permissions: %v", err)
	}

	// Run the install script with --method standalone and --prefix.
	cmd := exec.Command("sh", tempScript, "--method", "standalone", "--prefix", codeServerInstallDir)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Failed to run install script: %v", err)
	}

	fmt.Print(codeServerBinaryPath)

	// Verify the binary exists by checking the symlink and its target.
	if _, err := os.Lstat(codeServerBinaryPath); err != nil {
		return fmt.Sprintf("Failed to install code-server: symlink not found at %s", codeServerBinaryPath)
	}
	target, err := os.Readlink(codeServerBinaryPath)
	if err != nil {
		return fmt.Sprintf("Failed to read symlink %s: %v", codeServerBinaryPath, err)
	}
	if _, err := os.Stat(target); err != nil {
		return fmt.Sprintf("Failed to install code-server: binary target not found at %s", target)
	}
	// Clean up the install script.
	os.Remove(tempScript)

	return "Successfully installed Code Server"
}

// StartCodeServer launches code-server bound to a Unix socket.
// Uses a config file (config.yaml) for settings and runs as a subprocess with minimal environment.
func StartCodeServer() error {
	// Verify the code-server binary symlink and its target exist.
	if _, err := os.Lstat(codeServerBinaryPath); err != nil {
		logger.Main.Error("Code-server binary symlink not found: " + err.Error())
		return fmt.Errorf("code-server binary symlink not found at %s: %v", codeServerBinaryPath, err)
	}
	target, err := os.Readlink(codeServerBinaryPath)
	if err != nil {
		logger.Main.Error("Failed to read code-server binary symlink: " + err.Error())
		return fmt.Errorf("failed to read code-server binary symlink %s: %v", codeServerBinaryPath, err)
	}
	if _, err := os.Stat(target); err != nil {
		logger.Main.Error("Code-server binary target not found: " + err.Error())
		return fmt.Errorf("code-server binary target not found at %s: %v", target, err)
	}
	logger.Main.Info(fmt.Sprintf("Resolved code-server binary symlink %s to %s", codeServerBinaryPath, target))

	// Create config file at config.yaml with minimal settings.
	configContent := `auth: none
disable-telemetry: true
disable-update-check: true
disable-workspace-trust: true
disable-file-uploads: true
disable-getting-started-override: true
ignore-last-opened: true
`

	settingsContent := `{
    "workbench.colorTheme": "Solarized Dark"
}
`

	// Ensure directories exist.
	if err := os.MkdirAll(filepath.Dir(settingsFilePath), 0755); err != nil {
		logger.Main.Error("Failed to create CodeServer settings directory: " + err.Error())
		return fmt.Errorf("failed to create CodeServer settings directory: %v", err)
	}
	if err := os.MkdirAll(codeServerUserDataDir, 0755); err != nil {
		logger.Main.Error("Failed to create CodeServer user data directory: " + err.Error())
		return fmt.Errorf("failed to create CodeServer user data directory: %v", err)
	}
	if err := os.MkdirAll(codeServerExtensionsDir, 0755); err != nil {
		logger.Main.Error("Failed to create CodeServer extensions directory: " + err.Error())
		return fmt.Errorf("failed to create CodeServer extensions directory: %v", err)
	}

	// Write config.yaml.
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		logger.Main.Error("Failed to create CodeServer config file: " + err.Error())
		return fmt.Errorf("failed to create CodeServer config file: %v", err)
	}

	// Write settings.json, but only if it doesn't exist.
	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(settingsFilePath, []byte(settingsContent), 0644); err != nil {
			logger.Main.Error("Failed to create CodeServer settings file: " + err.Error())
			return fmt.Errorf("failed to create CodeServer settings file: %v", err)
		}
	}

	// Log the command we're about to run.
	logger.Main.Info(fmt.Sprintf("Starting code-server from %s", target))

	cmd := exec.Command(
		target, // Use the resolved target path instead of the symlink
		"--socket", codeServerSocketPath,
		"--socket-mode", "600",
		"--config", configFilePath,
		"--user-data-dir", codeServerUserDataDir,
		"--extensions-dir", codeServerExtensionsDir,
		//"--verbose",
	)

	if cmd.Err != nil {
		logger.Main.Error("Failed to create code-server command: " + cmd.Err.Error())
		return fmt.Errorf("failed to create code-server command: %v", cmd.Err)
	}

	// Set minimal environment
	// Set minimal environment variables, including PATH with code-server's bin directory.
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + filepath.Join(codeServerInstallDir, "bin") + ":" + os.Getenv("PATH"),
	}

	// (uncomment for debugging).
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	// Start the process.
	if err := cmd.Start(); err != nil {
		logger.Main.Error("Failed to start code-server: " + err.Error())
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	go func() {
		time.Sleep(600 * time.Millisecond)
		// Check if the socket exists to confirm code-server is running.
		if _, err := os.Stat(codeServerSocketPath); os.IsNotExist(err) {
			logger.Main.Warn("Expected Code-server socket was not found after 600ms: " + err.Error())
			return
		}
	}()

	go func() {
		if err := cmd.Wait(); err != nil {
			logger.Main.Error("Code-server exited with error: " + err.Error())
			fmt.Printf("code-server exited with error: %v\n", err)
		}
	}()

	return nil
}
