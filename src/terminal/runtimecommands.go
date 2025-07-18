// Package misc provides a non-blocking command-line interface for entering commands
// while allowing the application to continue its operations normally.
package terminal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/loader"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/steammgr"
)

// ANSI escape codes for green text and reset
const (
	cliPrompt = "\033[32m" + "SSUICLI" + " » " + "\033[0m"
)

// CommandFunc defines the signature for command handler functions.
type CommandFunc func(args []string) error

// commandRegistry holds the map of command names to their handler functions.
var commandRegistry = make(map[string]CommandFunc)
var mu sync.Mutex

var commandAliases = make(map[string][]string)

// RegisterCommand adds a new command and its handler to the registry.
func RegisterCommand(name string, handler CommandFunc, aliases ...string) {
	mu.Lock()
	defer mu.Unlock()
	commandRegistry[name] = handler
	if len(aliases) > 0 {
		commandAliases[name] = append(commandAliases[name], aliases...)
		for _, alias := range aliases {
			commandRegistry[alias] = handler
		}
	}
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
		time.Sleep(10 * time.Millisecond)

		for {
			fmt.Print(cliPrompt)
			os.Stdout.Sync() // Force flush the output buffer
			if !scanner.Scan() {
				break
			}
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

// helpCommand displays available commands along with their aliases.
func helpCommand(args []string) error {
	mu.Lock()
	defer mu.Unlock()
	logger.Core.Info("Available commands:")
	// Collect primary commands (those in commandAliases keys)
	primaryCommands := make([]string, 0, len(commandAliases))
	for cmd := range commandAliases {
		primaryCommands = append(primaryCommands, cmd)
	}
	sort.Strings(primaryCommands)
	for _, cmd := range primaryCommands {
		aliases := commandAliases[cmd]
		if len(aliases) > 0 {
			logger.Core.Infof("- %s (aliases: %s)", cmd, strings.Join(aliases, ", "))
		} else {
			logger.Core.Infof("- %s", cmd)
		}
	}
	return nil
}

// init registers default cli commands and their aliases.
func init() {
	RegisterCommand("help", helpCommand, "h")
	RegisterCommand("reloadbackend", WrapNoReturn(loader.ReloadAll), "rlb", "rb", "r")
	RegisterCommand("reloadconfig", WrapNoReturn(loader.ReloadConfig), "rlc", "rc")
	RegisterCommand("restartbackend", WrapNoReturn(loader.RestartBackend), "rsb")
	RegisterCommand("runsteamcmd", WrapNoReturn(steammgr.RunSteamCMD), "runsteam", "st")
	RegisterCommand("exit", WrapNoReturn(exitfromcli), "e")
	RegisterCommand("deleteconfig", WrapNoReturn(deleteConfig), "delc", "dc")
	RegisterCommand("sendtelemetry", WrapNoReturn(inop), "sendtel", "tel")
}

func exitfromcli() {
	// send signal to the main process to exit
	logger.Core.Info("I have to go...")
	os.Exit(0)
}

func inop() {
	logger.Core.Info("Not implemented yet")
}

func deleteConfig() {
	//remove file at config.ConfigPath
	if err := os.Remove(config.ConfigPath); err != nil {
		logger.Core.Error("Error deleting config file: " + err.Error())
		return
	}
	logger.Core.Info("Config file deleted successfully")
}
