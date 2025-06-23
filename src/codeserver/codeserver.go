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
)

// Paths defined at the top for clarity and maintainability.
const (
	// CodeServerBinaryPath is where the code-server binary will be installed.
	CodeServerBinaryPath = "/usr/bin/code-server"
	// CodeServerSocketPath is the Unix socket for internal-only communication.
	CodeServerSocketPath = "./cs/codeserver.sock"
	// GameServerDir is the directory containing game server files.
	GameServerDir = "./"
	// InstallScriptURL is the official code-server install script.
	InstallScriptURL = "https://code-server.dev/install.sh"
	// ConfigFilePath is the code-server configuration file.
	ConfigFilePath = "./cs/config.yaml"
)

// InitCodeServer initializes code-server at server startup.
// Creates the ./cs directory, installs, and starts code-server.
func InitCodeServer() error {

	os.RemoveAll("./cs")
	// Create ./cs directory if it doesn't exist.
	if err := os.MkdirAll("./cs", 0755); err != nil {
		return fmt.Errorf("failed to create cs directory: %v", err)
	}

	fmt.Print("Initializing Code Server...")

	// Install code-server.
	fmt.Println("Downloading Code Server")
	msg := DownloadInstallCodeServer()
	fmt.Println(msg)
	if !strings.Contains(strings.ToLower(msg), "successfully") && !strings.Contains(strings.ToLower(msg), "already installed") {
		return fmt.Errorf("code-server installation failed: %s", msg)
	}

	fmt.Print("Starting Code Server...")
	// Start code-server.
	if err := StartCodeServer(); err != nil {
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	return nil
}

// DownloadInstallCodeServer downloads and installs code-server using the official install script.
// Installs to ./cs/bin/code-server, keeping everything in the current directory.
// Returns a string with the result (success or error message).
func DownloadInstallCodeServer() string {
	// Enforce Linux-only support.
	if strings.ToLower(runtime.GOOS) != "linux" {
		return "Code Server is only supported on Linux"
	}

	// Check if code-server binary already exists to avoid re-installing.
	if _, err := os.Stat(CodeServerBinaryPath); err == nil {
		return "Code Server already installed"
	}

	// Create ./cs/bin directory for the binary.
	binDir := filepath.Dir(CodeServerBinaryPath)
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Sprintf("Failed to create bin directory: %v", err)
	}

	// Create a temporary file for the install script.
	tempScript := "./cs/install.sh"
	if err := os.MkdirAll(filepath.Dir(tempScript), 0755); err != nil {
		return fmt.Sprintf("Failed to create cs directory: %v", err)
	}

	// Download the install script.
	resp, err := http.Get(InstallScriptURL)
	if err != nil {
		return fmt.Sprintf("Failed to download install script: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Failed to download install script: HTTP %d", resp.StatusCode)
	}

	// Save the script to ./cs/install.sh.
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
	cmd.Env = append(os.Environ(), "DESTDIR=./cs", "PREFIX=")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Failed to run install script: %v", err)
	}

	// Verify the binary exists.
	if _, err := os.Stat(CodeServerBinaryPath); os.IsNotExist(err) {
		return "Failed to install code-server: binary not found"
	}

	// Set executable permissions for the binary.
	if err := os.Chmod(CodeServerBinaryPath, 0755); err != nil {
		return fmt.Sprintf("Failed to set binary permissions: %v", err)
	}

	// Clean up the install script.
	os.Remove(tempScript)

	return "Successfully installed Code Server"
}

// StartCodeServer launches code-server bound to a Unix socket, restricted to GameServerDir.
// Uses a config file (./cs/config.yaml) for settings and runs as a subprocess with minimal environment.
func StartCodeServer() error {

	// Create config file at ./cs/config.yaml with minimal settings.
	configContent := `auth: none
disable-telemetry: true
disable-update-check: true
disable-workspace-trust: true
disable-file-uploads: true
user-data-dir: ./cs/user-data
extensions-dir: ./cs/extensions
disable-getting-started-override: true
ignore-last-opened: true

`
	if err := os.WriteFile(ConfigFilePath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}

	cmd := exec.Command(
		CodeServerBinaryPath,
		"--socket", CodeServerSocketPath,
		"--socket-mode", "600",
		"--config", ConfigFilePath,
		//"--verbose",
	)

	// Set minimal environment variables (HOME and PATH).
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
	}

	// Capture stdout/stderr for verbose logging.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	// Wait briefly to check if the socket is created.
	time.Sleep(3 * time.Second)

	// Check if the socket exists to confirm code-server is running.
	if _, err := os.Stat(CodeServerSocketPath); os.IsNotExist(err) {
		return fmt.Errorf("code-server did not create socket: %v", err)
	}

	// Run in background (don't wait for cmd.Wait, as it's a long-running process).
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("code-server exited with error: %v\n", err)
		}
	}()

	return nil
}
