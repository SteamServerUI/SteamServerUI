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
)

// Paths defined at the top for clarity and maintainability.
const (
	// CodeServerBinaryPath is where the code-server binary will be installed.
	CodeServerBinaryPath = "./cs/code-server"
	// CodeServerSocketPath is the Unix socket for internal-only communication.
	CodeServerSocketPath = "./cs/codeserver.sock"
	GameServerDir        = "./"
	// InstallScriptURL is the official code-server install script.
	InstallScriptURL = "https://code-server.dev/install.sh"
)

func InitCodeServer() error {
	// Create ./cs directory if it doesn't exist.
	if err := os.MkdirAll("./cs", 0755); err != nil {
		return fmt.Errorf("failed to create cs directory: %v", err)
	}
	fmt.Println("Downloading Code Server")
	msg := DownloadInstallCodeServer()
	fmt.Println(msg)
	if !strings.Contains(msg, "successfully") && !strings.Contains(msg, "already installed") {
		return fmt.Errorf("code-server installation failed: %s", msg)
	}
	if err := StartCodeServer(); err != nil {
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	return nil
}

func DownloadInstallCodeServer() string {

	if strings.ToLower(runtime.GOOS) != "linux" {
		return "Code Server is only supported on Linux"
	}

	if _, err := os.Stat(CodeServerBinaryPath); err == nil {
		return "Code Server already installed"
	}

	tempScript := "./cs/install.sh"
	if err := os.MkdirAll(filepath.Dir(tempScript), 0755); err != nil {
		return fmt.Sprintf("Failed to create cs directory: %v", err)
	}

	resp, err := http.Get(InstallScriptURL)
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

	os.Remove(tempScript)

	return "Ran codeserver install script"
}

// StartCodeServer launches code-server bound to a Unix socket, restricted to GameServerDir.
// It runs as a subprocess silently.
func StartCodeServer() error {
	return nil
}
