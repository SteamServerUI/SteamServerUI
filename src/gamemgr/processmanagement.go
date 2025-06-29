package gamemgr

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v6/src/argmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/ssestream"
	"github.com/google/uuid"
)

var (
	cmd     *exec.Cmd
	mu      sync.Mutex
	err     error
	exePath string
	logDone chan struct{} // Remove initialization here
)

// InternalIsServerRunning checks if the server process is running.
// Safe to call standalone as it manages its own locking.
func InternalIsServerRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return internalIsServerRunningNoLock()
}

// InternalStartServer starts the game server process.
func InternalStartServer() error {
	mu.Lock()
	defer mu.Unlock()

	if internalIsServerRunningNoLock() {
		return fmt.Errorf("server is already running")
	}

	// Initialize logDone for this server instance
	logDone = make(chan struct{})

	var args []string
	var err error
	args, err = argmgr.BuildCommandArgs()
	if err != nil {
		return fmt.Errorf("failed to build command args: %v", err)
	}

	exePath, err = getExePath()
	if err != nil {
		return fmt.Errorf("failed to get exePath: %w", err)
	}

	logger.Core.Info("=== GAMESERVER STARTING ===")
	logger.Core.Info("• Executable: " + exePath)
	logger.Core.Info("• Arguments: " + strings.Join(args, " "))

	// Delegate to platform-specific start logic
	if err := platformStartServer(exePath, args); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	// Create a UUID for this specific run
	createGameServerUUID()
	logger.Core.Debug("Created Game Server with internal UUID: " + config.GetGameServerUUID().String())
	return nil
}

// InternalStopServer stops the game server process.
func InternalStopServer() error {
	mu.Lock()
	defer mu.Unlock()

	if !internalIsServerRunningNoLock() {
		return fmt.Errorf("server not running")
	}

	// Delegate to platform-specific stop logic
	if err := platformStopServer(); err != nil {
		return fmt.Errorf("error stopping server: %v", err)
	}

	// Process is confirmed stopped, clear cmd
	cmd = nil
	clearGameServerUUID()
	return nil
}

// internalIsServerRunningNoLock checks if the server process is running.
// Caller must hold mu.Lock().
func internalIsServerRunningNoLock() bool {
	// Implementation moved to platform-specific files
	return platformIsServerRunningNoLock()
}

// clearGameServerUUID clears the game server UUID.
func clearGameServerUUID() {
	config.SetGameServerUUID(uuid.Nil)
}

// createGameServerUUID creates a new game server UUID.
func createGameServerUUID() {
	config.SetGameServerUUID(uuid.New())
}

// getExePath returns the exePath based on argmgr.CurrentRunfile and runtime.GOOS.
func getExePath() (string, error) {
	if runtime.GOOS == "windows" {
		return argmgr.CurrentRunfile.WindowsExecutable, nil
	}
	if runtime.GOOS == "linux" {
		return argmgr.CurrentRunfile.LinuxExecutable, nil
	}
	return "", fmt.Errorf("no executable path found")
}

// readPipe reads from a pipe and broadcasts the output.
func readPipe(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	logger.Core.Debug("Started reading pipe")
	for scanner.Scan() {
		output := scanner.Text()
		ssestream.BroadcastConsoleOutput(output)
	}
	if err := scanner.Err(); err != nil {
		logger.Core.Debug("Pipe error: " + err.Error())
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading pipe: %v", err))
	}
	logger.Core.Debug("Pipe closed")
}
