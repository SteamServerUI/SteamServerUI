package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

var (
	runningPlugins      map[string]*exec.Cmd
	runningPluginsMutex sync.Mutex
	processExits        map[string]chan struct{}
)

// Initialize the runningPlugins and processExits maps
func init() {
	runningPlugins = make(map[string]*exec.Cmd)
	processExits = make(map[string]chan struct{})
}

// isPluginRunning checks if a plugin process is running
func isPluginRunning(pluginname string) bool {
	runningPluginsMutex.Lock()
	defer runningPluginsMutex.Unlock()

	cmd, exists := runningPlugins[pluginname]
	if !exists || cmd == nil || cmd.Process == nil {
		return false
	}

	if runtime.GOOS == "windows" {
		select {
		case <-processExits[pluginname]:
			delete(runningPlugins, pluginname)
			delete(processExits, pluginname)
			return false
		default:
			return true
		}
	}

	// On Linux, use Signal(0) to check if process is alive
	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		logger.Plugin.Debug(fmt.Sprintf("Signal(0) failed for plugin %s, assuming dead: %v", pluginname, err))
		delete(runningPlugins, pluginname)
		delete(processExits, pluginname)
		return false
	}
	return true
}

// StopPlugin stops a running plugin
func StopPlugin(pluginname string) error {
	runningPluginsMutex.Lock()
	defer runningPluginsMutex.Unlock()

	cmd, exists := runningPlugins[pluginname]
	if !exists || cmd == nil || cmd.Process == nil {
		return fmt.Errorf("plugin %s not running", pluginname)
	}

	isWindows := runtime.GOOS == "windows"
	var killErr error

	if isWindows {
		// On Windows, terminate the process
		killErr = cmd.Process.Kill()
		if processExits[pluginname] != nil {
			select {
			case <-processExits[pluginname]:
				logger.Plugin.Debug(fmt.Sprintf("processExits channel confirmed plugin %s shutdown", pluginname))
			case <-time.After(2 * time.Second):
				logger.Plugin.Warn(fmt.Sprintf("Timeout waiting for processExits confirmation for plugin %s", pluginname))
			}
		}
	} else {
		// On Linux, send SIGTERM for graceful shutdown
		if termErr := cmd.Process.Signal(syscall.SIGTERM); termErr != nil {
			logger.Plugin.Debug(fmt.Sprintf("SIGTERM failed for plugin %s: %v", pluginname, termErr))
			killErr = cmd.Process.Kill()
		} else {
			waitErrChan := make(chan error, 1)
			go func() {
				waitErrChan <- cmd.Wait()
			}()

			select {
			case waitErr := <-waitErrChan:
				if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
					logger.Plugin.Debug(fmt.Sprintf("Wait error after SIGTERM for plugin %s: %v", pluginname, waitErr))
				}
			case <-time.After(10 * time.Second):
				logger.Plugin.Warn(fmt.Sprintf("Timeout waiting for graceful shutdown of plugin %s, sending SIGKILL", pluginname))
				killErr = cmd.Process.Kill()
				select {
				case waitErr := <-waitErrChan:
					if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
						logger.Plugin.Debug(fmt.Sprintf("Wait error after SIGKILL for plugin %s: %v", pluginname, waitErr))
					}
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
	delete(runningPlugins, pluginname)
	delete(processExits, pluginname)
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

			runningPluginsMutex.Lock()
			runningPlugins[pluginname] = cmd
			runningPluginsMutex.Unlock()

			if runtime.GOOS == "windows" {

				if err := cmd.Start(); err != nil {
					logger.Plugin.Error(fmt.Sprintf("Error starting plugin %s: %v", pluginname, err))
					continue
				}
				logger.Plugin.Info(fmt.Sprintf("Plugin %s started with PID: %d", pluginname, cmd.Process.Pid))

				// Monito exit
				processExits[pluginname] = make(chan struct{})
				go func(pname string) {
					err := cmd.Wait()
					if err != nil {
						logger.Plugin.Error(fmt.Sprintf("Plugin %s exited with error: %v", pname, err))
					} else {
						logger.Plugin.Warn(fmt.Sprintf("Plugin %s exited unexpectedly without error", pname))
					}
					runningPluginsMutex.Lock()
					delete(runningPlugins, pname)
					delete(processExits, pname)
					runningPluginsMutex.Unlock()
					if err := UnregisterPlugin(pname); err != nil {
						logger.Plugin.Error(fmt.Sprintf("Failed to unregister plugin %s after exit: %v", pname, err))
					}
				}(pluginname)
			} else {

				if err := cmd.Start(); err != nil {
					logger.Plugin.Error(fmt.Sprintf("Error starting plugin %s: %v", pluginname, err))
					continue
				}
				logger.Plugin.Info(fmt.Sprintf("Plugin %s started with PID: %d", pluginname, cmd.Process.Pid))

				// Monitor process exit
				processExits[pluginname] = make(chan struct{})
				go func(pname string) {
					err := cmd.Wait()
					if err != nil {
						if err.Error() != "signal: terminated" {
							logger.Plugin.Debug(fmt.Sprintf("Plugin %s exited with error: %v", pname, err))
						}
					} else {
						logger.Plugin.Warn(fmt.Sprintf("Plugin %s exited unexpectedly without error", pname))
					}
					runningPluginsMutex.Lock()
					delete(runningPlugins, pname)
					delete(processExits, pname)
					runningPluginsMutex.Unlock()
					if err := UnregisterPlugin(pname); err != nil {
						logger.Plugin.Error(fmt.Sprintf("Failed to unregister plugin %s after exit: %v", pname, err))
					}
				}(pluginname)
				logger.Plugin.Warn("End of goroutine")
			}
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
