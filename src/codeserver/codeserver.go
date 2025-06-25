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

// Paths defined at the top for clarity and maintainability.
var (
	codeServerPath       = config.CodeServerPath
	codeServerBinaryPath = config.CodeServerBinaryPath
	codeServerSocketPath = config.CodeServerSocketPath
	installScriptURL     = config.InstallScriptURL
	configFilePath       = config.ConfigFilePath
)

// InitCodeServer initializes code-server at server startup.
// Creates the directory, installs, and starts code-server.
func InitCodeServer() error {

	if !config.GetIsCodeServerEnabled() {
		return nil
	}

	os.RemoveAll(codeServerPath)
	// Create directory if it doesn't exist.
	if err := os.MkdirAll(codeServerPath, 0755); err != nil {
		return fmt.Errorf("failed to create cs directory: %v", err)
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
	if _, err := os.Stat(codeServerBinaryPath); err == nil {
		return "Code Server already installed"
	}

	// Create a temporary file for the install script.
	tempScript := codeServerPath + "/install.sh"
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

	cmd := exec.Command("sh", tempScript)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Failed to run install script: %v", err)
	}

	// Verify the binary exists.
	if _, err := os.Stat(codeServerBinaryPath); os.IsNotExist(err) {
		return "Failed to install code-server: binary not found"
	}

	// Clean up the install script.
	os.Remove(tempScript)

	return "Successfully installed Code Server"
}

// StartCodeServer launches code-server bound to a Unix socket.
// Uses a config file (config.yaml) for settings and runs as a subprocess with minimal environment.
func StartCodeServer() error {

	// Create config file at config.yaml with minimal settings.
	configContent := `auth: none
disable-telemetry: true
disable-update-check: true
disable-workspace-trust: true
disable-file-uploads: true
disable-getting-started-override: true
ignore-last-opened: true

`
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}

	cmd := exec.Command(
		codeServerBinaryPath,
		"--socket", codeServerSocketPath,
		"--socket-mode", "600",
		"--config", configFilePath,
		//"--verbose",
	)

	// Set minimal environment variables (HOME and PATH).
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
	}

	// Capture stdout/stderr for verbose logging.
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	// Start the process.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	// Wait briefly to check if the socket is created.
	time.Sleep(3 * time.Second)

	// Check if the socket exists to confirm code-server is running.
	if _, err := os.Stat(codeServerSocketPath); os.IsNotExist(err) {
		return fmt.Errorf("code-server did not create socket: %v", err)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("code-server exited with error: %v\n", err)
		}
	}()

	return nil
}
