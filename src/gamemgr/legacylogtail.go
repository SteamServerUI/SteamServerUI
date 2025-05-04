//go:build linux

package gamemgr

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/ssestream"
)

// tailLogFile uses tail to read the log file because using the gameserver's output in pipes to read the serverlog doesn't work on Linux with the Stationeers gameserver.
// I didn't manage to implement proper file tailing (tail behavior) here in go, so I opted to just use the actual tail.. This is a workaround for a workaround.

func tailLogFile(logFilePath string) {
	logger.Core.Warn("LegacyLogFile tail of " + logFilePath + " enabled")

	// Wait and retry until the log file exists
	for i := range 10 { // Retry up to 10 times
		if _, err := os.Stat(logFilePath); err == nil {
			break // File exists, proceed
		}
		logger.Core.Debug("Log file " + logFilePath + " not found, retrying in 1s (" + strconv.Itoa(i+1) + "/10)")
		time.Sleep(1 * time.Second)
	}

	// If file still doesn't exist, give up and report
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		logger.Core.Debug("Log file " + logFilePath + " still not found after retries")
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Log file %s not found after retries", logFilePath))
		return
	}

	// Start tail -F (robust against rotation)
	cmd := exec.Command("tail", "-F", logFilePath)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Core.Debug("Error creating stdout pipe for tail: " + err.Error())
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error starting tail -F: %v", err))
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		logger.Core.Debug("Error starting tail -F: " + err.Error())
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error starting tail -F: %v", err))
		return
	}

	// Channel to signal when tail process should stop
	stopTail := make(chan struct{})
	defer func() {
		close(stopTail) // Signal goroutine to stop
		if err := cmd.Process.Kill(); err != nil {
			logger.Core.Debug("Error killing tail process: " + err.Error())
		}
		if err := cmd.Wait(); err != nil {
			logger.Core.Debug("Tail process exited with: " + err.Error())
		}
	}()

	scanner := bufio.NewScanner(pipe)
	logger.Core.Debug("Started tailing log file with tail -F")

	// Goroutine to read and broadcast tail output
	go func() {
		defer pipe.Close() // Close pipe when goroutine exits
		for scanner.Scan() {
			select {
			case <-stopTail:
				return // Exit if stop signal received
			default:
				output := scanner.Text()
				ssestream.BroadcastConsoleOutput(output)
			}
		}
		if err := scanner.Err(); err != nil {
			logger.Core.Debug("Error reading tail -F output: " + err.Error())
			ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading tail -F output: %v", err))
		}
	}()

	// Wait for logDone signal or timeout
	select {
	case <-logDone:
		logger.Core.Debug("Received logDone signal, stopping tail -F")
	case <-time.After(24 * time.Hour): // Arbitrary long timeout to prevent indefinite hang
		logger.Core.Warn("Timeout waiting for logDone signal, stopping tail -F")
	}
}
