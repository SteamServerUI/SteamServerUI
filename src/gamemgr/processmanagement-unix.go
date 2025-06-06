//go:build linux

package gamemgr

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

func platformIsServerRunningNoLock() bool {
	if cmd == nil || cmd.Process == nil {
		return false
	}

	// Use non-blocking check to avoid hanging
	_, err := syscall.Wait4(cmd.Process.Pid, nil, syscall.WNOHANG, nil)
	if err != nil {
		logger.Core.Debug("Wait4 failed, assuming process is dead: " + err.Error())
		cmd = nil
		clearGameServerUUID()
		return false
	}
	return true
}

func platformStartServer(exePath string, args []string) error {
	var envVars []string
	var err error

	if config.GetIsSSCMEnabled() {
		envVars, err = SetupBepInExEnvironment()
		if err != nil {
			return err
		}
	}
	cmd = exec.Command(exePath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // Create new process group
	if envVars != nil {
		cmd.Env = envVars
		logger.Core.Info("BepInEx/Doorstop environment configured for server process")
	}

	isLegacyLogMode := config.GetLegacyLogFile() != ""
	if isLegacyLogMode {
		go tailLogFile(config.GetLegacyLogFile()) // Run in goroutine to avoid blocking
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	logger.Core.Debug("Created pipes")

	go readPipe(stdout)
	go readPipe(stderr)
	return nil
}

func platformStopServer() error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	logger.Core.Debug("Stopping server with PID: " + strconv.Itoa(cmd.Process.Pid))

	// Send SIGTERM to the process group
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		logger.Core.Debug("Failed to send SIGTERM to process group: " + err.Error())
		return cmd.Process.Kill() // Fallback to killing the main process
	}

	waitErrChan := make(chan error, 1)
	go func() {
		waitErrChan <- cmd.Wait()
	}()

	select {
	case waitErr := <-waitErrChan:
		if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
			logger.Core.Debug("Wait error after SIGTERM: " + waitErr.Error())
		}
	case <-time.After(10 * time.Second):
		logger.Core.Warn("Timeout waiting for graceful shutdown, sending SIGKILL to process group")
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			logger.Core.Warn("Failed to send SIGKILL to process group: " + err.Error())
			return cmd.Process.Kill()
		}
		select {
		case waitErr := <-waitErrChan:
			if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
				logger.Core.Debug("Wait error after SIGKILL: " + waitErr.Error())
			}
		case <-time.After(2 * time.Second):
			return fmt.Errorf("timeout waiting for process to exit after SIGKILL")
		}
	}

	// Signal that the server has stopped
	if logDone != nil {
		close(logDone)
	}
	logger.Core.Debug("Server stopped, logDone signal sent")
	return nil
}
