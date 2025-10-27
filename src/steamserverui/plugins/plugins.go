package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

var (
	RunningPlugins      map[string]*exec.Cmd
	RunningPluginsMutex sync.Mutex
	processExits        map[string]chan struct{}
	// Signal intentional stops to prevent auto-unregister
	intentionalStops map[string]chan struct{}
)

// Initialize the maps
func init() {
	RunningPlugins = make(map[string]*exec.Cmd)
	processExits = make(map[string]chan struct{})
	intentionalStops = make(map[string]chan struct{})
}

// isPluginRunning checks if a plugin process is running
func isPluginRunning(pluginname string) bool {
	RunningPluginsMutex.Lock()
	defer RunningPluginsMutex.Unlock()

	cmd, exists := RunningPlugins[pluginname]
	if !exists || cmd == nil || cmd.Process == nil {
		return false
	}

	if runtime.GOOS == "windows" {
		select {
		case <-processExits[pluginname]:
			delete(RunningPlugins, pluginname)
			delete(processExits, pluginname)
			delete(intentionalStops, pluginname)
			return false
		default:
			return true
		}
	}

	// On Linux, use Signal(0) to check if process is alive
	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		logger.Plugin.Debug(fmt.Sprintf("Signal(0) failed for plugin %s, assuming dead: %v", pluginname, err))
		delete(RunningPlugins, pluginname)
		delete(processExits, pluginname)
		delete(intentionalStops, pluginname)
		return false
	}
	return true
}

// StopPlugin stops a running plugin
func StopPlugin(pluginname string) error {
	RunningPluginsMutex.Lock()

	cmd, exists := RunningPlugins[pluginname]
	if !exists || cmd == nil || cmd.Process == nil {
		RunningPluginsMutex.Unlock()
		return fmt.Errorf("plugin %s not running", pluginname)
	}

	// Signal the monitor goroutine that this is an intentional stop
	stopChan := make(chan struct{})
	intentionalStops[pluginname] = stopChan
	close(stopChan)

	// Get the exit channel before unlocking
	exitChan := processExits[pluginname]

	RunningPluginsMutex.Unlock()

	isWindows := runtime.GOOS == "windows"
	var killErr error

	if isWindows {
		// On Windows, terminate the process
		killErr = cmd.Process.Kill()
		if exitChan != nil {
			select {
			case <-exitChan:
				logger.Plugin.Debug(fmt.Sprintf("processExits channel confirmed plugin %s shutdown", pluginname))
			case <-time.After(2 * time.Second):
				logger.Plugin.Warn(fmt.Sprintf("Timeout waiting for processExits confirmation for plugin %s", pluginname))
			}
		}
	} else {
		// On Linux, send SIGTERM and let the monitor goroutine handle Wait()
		if termErr := cmd.Process.Signal(syscall.SIGTERM); termErr != nil {
			logger.Plugin.Debug(fmt.Sprintf("SIGTERM failed for plugin %s: %v", pluginname, termErr))
			killErr = cmd.Process.Kill()
		}

		// Wait for the monitor goroutine to complete
		if exitChan != nil {
			select {
			case <-exitChan:
				logger.Plugin.Debug(fmt.Sprintf("Monitor goroutine confirmed plugin %s shutdown", pluginname))
			case <-time.After(10 * time.Second):
				logger.Plugin.Warn(fmt.Sprintf("Timeout waiting for graceful shutdown of plugin %s, sending SIGKILL", pluginname))
				killErr = cmd.Process.Kill()
				// Wait a bit more after SIGKILL
				select {
				case <-exitChan:
					logger.Plugin.Debug(fmt.Sprintf("Plugin %s stopped after SIGKILL", pluginname))
				case <-time.After(2 * time.Second):
					return fmt.Errorf("timeout waiting for plugin %s to exit after SIGKILL", pluginname)
				}
			}
		}
	}

	if killErr != nil {
		return fmt.Errorf("error stopping plugin %s: %v", pluginname, killErr)
	}

	// Clear plugin state
	RunningPluginsMutex.Lock()
	delete(RunningPlugins, pluginname)
	delete(processExits, pluginname)
	delete(intentionalStops, pluginname)
	RunningPluginsMutex.Unlock()

	return nil
}

// ManagePlugins starts or restarts all registered plugins
func ManagePlugins() error {
	registeredPlugins := config.GetRegisteredPlugins()

	go func() {
		// Check each registered plugin
		for pluginname, filename := range registeredPlugins {
			pluginPath := filepath.Join(config.GetPluginsFolder(), filename)

			// Verify file exists and is executable
			fileInfo, err := os.Stat(pluginPath)
			if os.IsNotExist(err) {
				logger.Plugin.Error(fmt.Sprintf("Plugin file %s does not exist", pluginPath))
				if err := UnregisterPlugin(pluginname); err != nil {
					logger.Plugin.Error(fmt.Sprintf("Failed to unregister plugin %s: %v", pluginname, err))
				}
				continue
			}
			if err != nil {
				logger.Plugin.Error(fmt.Sprintf("Failed to stat plugin file %s: %v", pluginPath, err))
				if err := UnregisterPlugin(pluginname); err != nil {
					logger.Plugin.Error(fmt.Sprintf("Failed to unregister plugin %s: %v", pluginname, err))
				}
				continue
			}
			if runtime.GOOS == "linux" && fileInfo.Mode().Perm()&0111 == 0 {
				logger.Plugin.Error(fmt.Sprintf("Plugin file %s is not executable", pluginPath))
				if err := UnregisterPlugin(pluginname); err != nil {
					logger.Plugin.Error(fmt.Sprintf("Failed to unregister plugin %s: %v", pluginname, err))
				}
				continue
			}

			// Check if plugin is running
			if isPluginRunning(pluginname) {
				// Stop running plugin for restart
				if err := StopPlugin(pluginname); err != nil {
					logger.Plugin.Error(fmt.Sprintf("Failed to stop plugin %s for restart: %v", pluginname, err))
					continue
				}
				logger.Plugin.Info(fmt.Sprintf("Stopped plugin %s for restart", pluginname))
			}

			// Start plugin
			logger.Plugin.Info(fmt.Sprintf("Starting plugin %s", pluginname))
			cmd := exec.Command(pluginPath)

			RunningPluginsMutex.Lock()
			RunningPlugins[pluginname] = cmd
			processExits[pluginname] = make(chan struct{})
			exitChan := processExits[pluginname]
			RunningPluginsMutex.Unlock()

			if err := cmd.Start(); err != nil {
				logger.Plugin.Error(fmt.Sprintf("Error starting plugin %s: %v", pluginname, err))
				RunningPluginsMutex.Lock()
				delete(RunningPlugins, pluginname)
				delete(processExits, pluginname)
				RunningPluginsMutex.Unlock()
				continue
			}

			logger.Plugin.Info(fmt.Sprintf("Plugin %s started with PID: %d", pluginname, cmd.Process.Pid))

			// Monitor process exit
			go func(pname string, pcmd *exec.Cmd, pExitChan chan struct{}) {
				defer close(pExitChan)

				err := pcmd.Wait()

				// Check if this was an intentional stop
				RunningPluginsMutex.Lock()
				stopChan, intentional := intentionalStops[pname]
				RunningPluginsMutex.Unlock()

				if intentional {
					select {
					case <-stopChan:
						// This was an intentional stop, don't unregister
						logger.Plugin.Debug(fmt.Sprintf("Plugin %s stopped intentionally", pname))
						RunningPluginsMutex.Lock()
						delete(intentionalStops, pname)
						RunningPluginsMutex.Unlock()
						return
					default:
					}
				}

				// This was an unexpected exit, log and unregister
				if err != nil {
					if runtime.GOOS == "windows" {
						logger.Plugin.Error(fmt.Sprintf("Plugin %s exited with error: %v", pname, err))
					} else {
						if err.Error() != "signal: terminated" {
							logger.Plugin.Error(fmt.Sprintf("Plugin %s exited with error: %v", pname, err))
						} else {
							logger.Plugin.Debug(fmt.Sprintf("Plugin %s exited with SIGTERM", pname))
						}
					}
				} else {
					logger.Plugin.Warn(fmt.Sprintf("Plugin %s exited unexpectedly without error", pname))
				}

				RunningPluginsMutex.Lock()
				delete(RunningPlugins, pname)
				delete(processExits, pname)
				delete(intentionalStops, pname)
				RunningPluginsMutex.Unlock()

				if err := UnregisterPlugin(pname); err != nil {
					logger.Plugin.Error(fmt.Sprintf("Failed to unregister plugin %s after exit: %v", pname, err))
				}
			}(pluginname, cmd, exitChan)
		}
	}()

	return nil
}

// UnregisterPlugin removes a plugin from the registeredPlugins config
func UnregisterPlugin(pluginname string) error {
	registered := config.GetRegisteredPlugins()
	if _, exists := registered[pluginname]; exists {
		delete(registered, pluginname)
		if err := config.SetRegisteredPlugins(registered); err != nil {
			logger.Plugin.Error(fmt.Sprintf("Failed to unregister plugin %s: %v", pluginname, err))
			return fmt.Errorf("failed to unregister plugin %s", pluginname)
		}
		logger.Plugin.Info(fmt.Sprintf("Plugin %s unregistered", pluginname))
	}
	return nil
}
