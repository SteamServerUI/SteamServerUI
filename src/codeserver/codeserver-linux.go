//go:build linux

package codeserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
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

// ProcessManager manages the code-server process lifecycle
type ProcessManager struct {
	mu      sync.RWMutex
	cmd     *exec.Cmd
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
}

var processManager = &ProcessManager{}

// InitCodeServer initializes code-server at server startup.
// Creates the directory, installs, and starts code-server.
func InitCodeServer() error {
	if !config.GetIsCodeServerEnabled() {
		return nil
	}

	// Enforce Linux-only support.
	if strings.ToLower(runtime.GOOS) != "linux" {
		return fmt.Errorf("code server can only be used on Linux")
	}

	// Set up graceful shutdown handling
	setupGracefulShutdown()

	// Ensure any existing instance is stopped
	if err := stopCodeServer(); err != nil {
		logger.Codeserver.Warn("Failed to stop existing code-server instance: " + err.Error())
	}

	// Clean up socket file
	os.RemoveAll(codeServerSocketPath)

	// Create directory if it doesn't exist.
	if err := os.MkdirAll(codeServerPath, 0755); err != nil {
		return fmt.Errorf("failed to create Code Server directory: %v", err)
	}

	logger.Codeserver.Info("Initializing Code Server...")
	msg := downloadInstallCodeServer()
	logger.Codeserver.Info(msg)
	if !strings.Contains(strings.ToLower(msg), "successfully") && !strings.Contains(strings.ToLower(msg), "already installed") {
		return fmt.Errorf("code-server installation failed: %s", msg)
	}

	logger.Codeserver.Debug("Starting Code Server...")
	// Start code-server.
	if err := startCodeServer(); err != nil {
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	return nil
}

// stopCodeServer gracefully stops the code-server process
func stopCodeServer() error {
	processManager.mu.Lock()
	defer processManager.mu.Unlock()

	if !processManager.running || processManager.cmd == nil || processManager.cmd.Process == nil {
		return nil
	}

	logger.Codeserver.Infof("Stopping code-server with PID %d...", processManager.cmd.Process.Pid)

	// Cancel the context to signal shutdown
	if processManager.cancel != nil {
		processManager.cancel()
	}

	// Check if process is still alive before sending signal
	if err := processManager.cmd.Process.Signal(syscall.Signal(0)); err != nil {
		// Process is already dead
		logger.Codeserver.Info("Code-server process already terminated")
		processManager.cmd = nil
		processManager.running = false
		os.RemoveAll(codeServerSocketPath)
		return nil
	}

	// Try graceful shutdown first (SIGTERM)
	if err := processManager.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		logger.Codeserver.Warn("Failed to send SIGTERM to code-server: " + err.Error())
		// Process might have already finished
		processManager.cmd = nil
		processManager.running = false
		os.RemoveAll(codeServerSocketPath)
		return nil
	}

	// Clean up
	processManager.cmd = nil
	processManager.running = false
	os.RemoveAll(codeServerSocketPath)

	return nil
}

// setupGracefulShutdown sets up signal handlers for graceful shutdown
func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c
		logger.Codeserver.Info("Received shutdown signal, stopping code-server...")
		if err := stopCodeServer(); err != nil {
			logger.Codeserver.Error("Failed to stop code-server during shutdown: " + err.Error())
		} else {
			logger.Codeserver.Info("Code-server shutdown completed successfully")
		}
		os.Exit(0)
	}()
}

// downloadInstallCodeServer downloads and installs code-server using the official install script.
func downloadInstallCodeServer() string {
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
	if config.GetLogLevel() == 10 {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Failed to run install script: %v", err)
	}

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

// startCodeServer launches code-server bound to a Unix socket.
// Uses a config file (config.yaml) for settings and runs as a subprocess with minimal environment.
func startCodeServer() error {
	processManager.mu.Lock()
	defer processManager.mu.Unlock()

	// Check if already running
	if processManager.running {
		return fmt.Errorf("code-server is already running")
	}

	// Verify the code-server binary symlink and its target exist.
	if _, err := os.Lstat(codeServerBinaryPath); err != nil {
		logger.Codeserver.Error("Code-server binary symlink not found: " + err.Error())
		return fmt.Errorf("code-server binary symlink not found at %s: %v", codeServerBinaryPath, err)
	}
	target, err := os.Readlink(codeServerBinaryPath)
	if err != nil {
		logger.Codeserver.Error("Failed to read code-server binary symlink: " + err.Error())
		return fmt.Errorf("failed to read code-server binary symlink %s: %v", codeServerBinaryPath, err)
	}
	if _, err := os.Stat(target); err != nil {
		logger.Codeserver.Error("Code-server binary target not found: " + err.Error())
		return fmt.Errorf("code-server binary target not found at %s: %v", target, err)
	}
	logger.Codeserver.Debug(fmt.Sprintf("Resolved code-server binary symlink %s to %s", codeServerBinaryPath, target))

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
		logger.Codeserver.Error("Failed to create CodeServer settings directory: " + err.Error())
		return fmt.Errorf("failed to create CodeServer settings directory: %v", err)
	}
	if err := os.MkdirAll(codeServerUserDataDir, 0755); err != nil {
		logger.Codeserver.Error("Failed to create CodeServer user data directory: " + err.Error())
		return fmt.Errorf("failed to create CodeServer user data directory: %v", err)
	}
	if err := os.MkdirAll(codeServerExtensionsDir, 0755); err != nil {
		logger.Codeserver.Error("Failed to create CodeServer extensions directory: " + err.Error())
		return fmt.Errorf("failed to create CodeServer extensions directory: %v", err)
	}

	// Write config.yaml.
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		logger.Codeserver.Error("Failed to create CodeServer config file: " + err.Error())
		return fmt.Errorf("failed to create CodeServer config file: %v", err)
	}

	// Write settings.json, but only if it doesn't exist.
	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(settingsFilePath, []byte(settingsContent), 0644); err != nil {
			logger.Codeserver.Error("Failed to create CodeServer settings file: " + err.Error())
			return fmt.Errorf("failed to create CodeServer settings file: %v", err)
		}
	}

	// Create context for process management
	processManager.ctx, processManager.cancel = context.WithCancel(context.Background())

	// Log the command we're about to run.
	logger.Codeserver.Debug(fmt.Sprintf("Starting code-server from %s", target))

	processManager.cmd = exec.CommandContext(
		processManager.ctx,
		target, // Use the resolved target path instead of the symlink
		"--socket", codeServerSocketPath,
		"--socket-mode", "600",
		"--config", configFilePath,
		"--user-data-dir", codeServerUserDataDir,
		"--extensions-dir", codeServerExtensionsDir,
		//"--verbose",
	)

	if processManager.cmd.Err != nil {
		logger.Codeserver.Error("Failed to create code-server command: " + processManager.cmd.Err.Error())
		return fmt.Errorf("failed to create code-server command: %v", processManager.cmd.Err)
	}

	// Set minimal environment variables, including PATH with code-server's bin directory.
	processManager.cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + filepath.Join(codeServerInstallDir, "bin") + ":" + os.Getenv("PATH"),
	}

	// Set process group to ensure proper cleanup of child processes
	processManager.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if config.GetLogLevel() == 10 {
		processManager.cmd.Stdout = os.Stdout
		processManager.cmd.Stderr = os.Stderr
	}

	// Start the process.
	if err := processManager.cmd.Start(); err != nil {
		logger.Codeserver.Error("Failed to start code-server: " + err.Error())
		processManager.cancel()
		return fmt.Errorf("failed to start code-server: %v", err)
	}

	processManager.running = true
	logger.Codeserver.Info(fmt.Sprintf("Code-server started with PID %d", processManager.cmd.Process.Pid))

	// Monitor the process in a separate goroutine
	go func() {
		// Check if the socket exists to confirm code-server is running.
		time.Sleep(600 * time.Millisecond)
		if _, err := os.Stat(codeServerSocketPath); os.IsNotExist(err) {
			logger.Codeserver.Warn("Expected Code-server socket was not found after 600ms: " + err.Error())
		}
	}()

	// Wait for process completion in a separate goroutine
	go func() {
		err := processManager.cmd.Wait()

		processManager.mu.Lock()
		defer processManager.mu.Unlock()

		processManager.running = false

		if err != nil {
			// Check if it was cancelled (expected shutdown)
			if processManager.ctx.Err() == context.Canceled {
				logger.Codeserver.Info("Code-server stopped as requested")
			}
		}

		// Clean up socket file
		os.RemoveAll(codeServerSocketPath)
	}()

	return nil
}
