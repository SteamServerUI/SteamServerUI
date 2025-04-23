//go:build windows

package gamemgr

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

func platformIsServerRunningNoLock() bool {
	if cmd == nil || cmd.Process == nil {
		return false
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			logger.Core.Debug("Wait failed: " + err.Error())
			if strings.Contains(err.Error(), "The handle is invalid") {
				cmd = nil
				clearGameServerUUID()
				return false
			}
		}
		cmd = nil
		clearGameServerUUID()
		return false
	case <-time.After(50 * time.Millisecond):
		return true
	}
}

func platformStartServer(exePath string, args []string) error {
	cmd = exec.Command(exePath, args...)

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
	logger.Core.Debug("Server process started with PID:" + string(cmd.Process.Pid))
	logger.Core.Debug("Created pipes")

	go readPipe(stdout)
	go readPipe(stderr)
	return nil
}

func platformStopServer() error {
	killErr := cmd.Process.Kill()
	waitErrChan := make(chan error, 1)
	go func() {
		waitErrChan <- cmd.Wait()
	}()

	select {
	case waitErr := <-waitErrChan:
		if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") &&
			!strings.Contains(waitErr.Error(), "The handle is invalid") {
			return waitErr
		}
	case <-time.After(1 * time.Second):
		return fmt.Errorf("timeout waiting for process to exit")
	}

	if killErr != nil {
		return killErr
	}
	return nil
}
