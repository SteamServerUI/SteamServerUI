// Package misc provides a non-blocking command-line interface for entering commands
// while allowing the application to continue its operations normally.
package misc

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/loader"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/steammgr"
)

// CommandFunc defines the signature for command handler functions.
type CommandFunc func(args []string) error

// commandRegistry holds the map of command names to their handler functions.
var commandRegistry = make(map[string]CommandFunc)
var mu sync.Mutex

// RegisterCommand adds a new command and its handler to the registry.
func RegisterCommand(name string, handler CommandFunc) {
	mu.Lock()
	defer mu.Unlock()
	commandRegistry[name] = handler
}

// StartConsole starts a non-blocking console input loop in a separate goroutine.
func StartConsole(wg *sync.WaitGroup) {

	if !config.GetIsConsoleEnabled() {
		logger.Core.Warn("Console is disabled, skipping...")
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		logger.Core.Info("Console input started. Type 'help' for commands.")

		for scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}
			ProcessCommand(input)
		}

		if err := scanner.Err(); err != nil {
			logger.Core.Errorf("Console input error: %v", err)
		}
		logger.Core.Info("Console input stopped.")
	}()
}

// ProcessCommand parses and executes a command from the input string.
func ProcessCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	commandName := strings.ToLower(args[0])
	args = args[1:] // Remove command name from args

	mu.Lock()
	handler, exists := commandRegistry[commandName]
	mu.Unlock()

	if !exists {
		logger.Core.Errorf("Unknown command: %s. Type 'help' for available commands.", commandName)
		return
	}

	if err := handler(args); err != nil {
		logger.Core.Errorf("Command '%s' failed: %v", commandName, err)
	}
}

// WrapNoReturn wraps a function with no return value to match CommandFunc.
func WrapNoReturn(fn func()) CommandFunc {
	return func(args []string) error {
		if len(args) > 0 {
			return errors.New("command does not accept arguments")
		}
		fn()
		logger.Core.Info("Runtime CLI Command executed successfully")
		return nil
	}
}

// helpCommand displays available commands.
func helpCommand(args []string) error {
	mu.Lock()
	defer mu.Unlock()
	logger.Core.Info("Available commands:")
	logger.Core.Info("Each command can be run with the alias: 'rl' or 'reload' for example.")
	for cmd := range commandRegistry {
		logger.Core.Infof("- %s", cmd)
	}
	return nil
}

// init registers default commands.
func init() {
	RegisterCommand("help", helpCommand)
	RegisterCommand("h", helpCommand)
	RegisterCommand("reload", WrapNoReturn(loader.ReloadAll))
	RegisterCommand("r", WrapNoReturn(loader.ReloadAll))
	RegisterCommand("reloadconfig", WrapNoReturn(loader.ReloadConfig))
	RegisterCommand("rc", WrapNoReturn(loader.ReloadConfig))
	RegisterCommand("restartbackend", WrapNoReturn(loader.RestartBackend))
	RegisterCommand("rsb", WrapNoReturn(loader.RestartBackend))
	RegisterCommand("getconfig", WrapNoReturn(loader.PrintConfigDetails))
	RegisterCommand("getc", WrapNoReturn(loader.PrintConfigDetails))
	RegisterCommand("steamcmd", WrapNoReturn(steammgr.RunSteamCMD))
	RegisterCommand("st", WrapNoReturn(steammgr.RunSteamCMD))

}
