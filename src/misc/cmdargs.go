package misc

import (
	"flag"
	"fmt"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/security"
)

// LoadCmdArgs parses command-line arguments ONCE at startup (called from func main) and applies them using the config setters.
// Because this is using the config rather than adding features to it, it is a part of the loader package.
func LoadCmdArgs() {
	// Define flags matching the config variable names
	var backendEndpointPort string
	var backendEndpointIP string
	var gameBranch string
	var logLevel int
	var isDebugMode bool
	var createSSUILogFile bool
	var recoveryPassword string
	var devMode string

	flag.StringVar(&backendEndpointPort, "BackendEndpointPort", "", "Override the backend endpoint port (e.g., 8080)")
	flag.StringVar(&backendEndpointPort, "port", "", "(Alias) Override the backend endpoint port (e.g., 8080)")
	flag.StringVar(&backendEndpointIP, "BackendEndpointIP", "", "Override the backend endpoint IP (e.g., 127.0.0.1)")
	flag.StringVar(&backendEndpointIP, "ip", "", "(Alias) Override the backend endpoint IP (e.g., 127.0.0.1)")
	flag.StringVar(&gameBranch, "GameBranch", "", "Override the game branch (e.g., beta)")
	flag.StringVar(&recoveryPassword, "RecoveryPassword", "", "Enable recovery user (expects password as argument)")
	flag.StringVar(&devMode, "dev", "", "This enables: Auth, OVERWRITES the admin user to admin:admin:superadmin, and enables the console. This is intended for development purposes only.")
	flag.IntVar(&logLevel, "LogLevel", 0, "Override the log level (e.g., 10)")
	flag.IntVar(&logLevel, "ll", 0, "(Alias) Override the log level (e.g., 10)")
	flag.BoolVar(&isDebugMode, "IsDebugMode", false, "Enable debug mode (true/false)")
	flag.BoolVar(&createSSUILogFile, "CreateSSUILogFile", false, "Create a log file for SSUI (true/false)")

	// Parse command-line flags
	flag.Parse()

	if devMode == "true" {
		config.SetAuthEnabled(true)
		config.SetIsFirstTimeSetup(false)
		config.SetUsers(map[string]string{"admin": "$2a$10$7QQhPkNAfT.MXhJhnnodXOyn3KKE/1eu7nYb0y2O1UBoAWc0Y/fda"}) // admin:admin
		config.SetUserLevels(map[string]string{"admin": "superadmin"})
		config.SetIsConsoleEnabled(true)
	}

	if backendEndpointPort != "" {
		oldPort := config.GetBackendEndpointPort()
		config.SetBackendEndpointPort(backendEndpointPort)
		logger.Main.Info(fmt.Sprintf("Overriding BackendEndpointPort from Command Line args: Before=%s, Now=%s", oldPort, backendEndpointPort))
	}

	if backendEndpointIP != "" {
		oldIP := config.GetBackendEndpointIP()
		config.SetBackendEndpointIP(backendEndpointIP)
		logger.Main.Info(fmt.Sprintf("Overriding BackendEndpointIP from Command Line args: Before=%s, Now=%s", oldIP, backendEndpointIP))
	}

	if gameBranch != "" {
		oldBranch := config.GetGameBranch()
		config.SetGameBranch(gameBranch)
		logger.Main.Info(fmt.Sprintf("Overriding GameBranch from Command Line args: Before=%s, Now=%s", oldBranch, gameBranch))
	}

	if recoveryPassword != "" {
		recoveryPassword := strings.TrimSpace(recoveryPassword)
		if recoveryPassword == "" {
			logger.Main.Error("IsRecoveryMode flag provided but password is empty. Skipping recovery user creation.")
		} else {
			hashedPassword, err := security.HashPassword(recoveryPassword)
			if err != nil {
				logger.Main.Error(fmt.Sprintf("Failed to hash recovery password: %v", err))
				return
			}
			config.SetUsers(map[string]string{"recovery": hashedPassword})
			config.SetUserLevels(map[string]string{"recovery": "superadmin"})
			logger.Main.Warn(fmt.Sprintf("Recovery user added with access level superadmin. Login with username 'recovery' and password '%s'", recoveryPassword))
		}
	}

	if logLevel != 0 {
		oldLevel := config.GetLogLevel()
		config.SetLogLevel(logLevel)
		logger.Main.Info(fmt.Sprintf("Overriding LogLevel from Command Line args: Before=%d, Now=%d", oldLevel, logLevel))
	}

	if flag.Lookup("IsDebugMode").Value.String() != "false" {
		oldDebug := config.GetIsDebugMode()
		config.SetIsDebugMode(isDebugMode)
		config.SetLogLevel(10)
		logger.Main.Info(fmt.Sprintf("Overriding IsDebugMode from Command Line args: Before=%t, Now=%t", oldDebug, isDebugMode))
	}
	if flag.Lookup("CreateSSUILogFile").Value.String() != "false" {
		oldCreateSSUILogFile := config.GetCreateSSUILogFile()
		config.SetCreateSSUILogFile(createSSUILogFile)
		logger.Main.Info(fmt.Sprintf("Overriding CreateSSUILogFile from Command Line args: Before=%t, Now=%t", oldCreateSSUILogFile, createSSUILogFile))
	}
}
