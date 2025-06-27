// loader.go
package loader

import (
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v6/src/argmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/backupmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/codeserver"
	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/discordbot"
	"github.com/SteamServerUI/SteamServerUI/v6/src/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/setup"
	"github.com/SteamServerUI/SteamServerUI/v6/src/steammgr"
)

func ReloadAll() {
	ReloadConfig()
	ReloadDiscordBot()
	ReloadRunfile()
	InitCodeServer()
	InitBackupMgr()
}

func ReloadConfig() {
	if _, err := config.LoadConfig(); err != nil {
		if !StillInSetup() {
			logger.Core.Error("Failed to load config: " + err.Error())
		}
		return
	}

	if StillInSetup() {
		logger.Core.Info("Config created or loaded successfully")
	} else {
		logger.Core.Info("Config reloaded successfully")
	}

	if config.GetIsSSCMEnabled() {
		setup.InstallBepInEx()
	}

	PrintConfigDetails()
}

func ReloadDiscordBot() {
	if config.GetIsDiscordEnabled() {
		go discordbot.InitializeDiscordBot()
		logger.Discord.Info("Discord bot reloaded successfully")
	}
}

func ReloadRunfile() error {
	if err := argmgr.LoadRunfile(config.GetRunfileGame(), config.GetRunFilesFolder()); err != nil {
		if !StillInSetup() {
			logger.Runfile.Error("Failed to reload runfile: " + err.Error())
		}
		return err
	}
	logger.Runfile.Info("Runfile reloaded successfully")
	return nil
}

func InitRunfile(game string) error {
	// Validate input
	game = strings.TrimSpace(game)
	if game == "" {
		return fmt.Errorf("game cannot be empty")
	}

	logger.Runfile.Info("Updating runfile game to " + game)
	logger.Runfile.Info("Stopping server if running")
	gamemgr.InternalStopServer()
	config.SetRunfileGame(game)

	if err := ReloadRunfile(); err != nil {
		return err
	}

	logger.Runfile.Info("Running SteamCMD, this may take a while...")
	steammgr.RunSteamCMD()
	logger.Runfile.Info("Runfile game updated to " + game)

	return nil
}

// The detector should NOT be reloaded, as it is a singleton. Instead, dynamic changes come in via the custom detections manager.
func InitDetector() {
	detector := detectionmgr.Start()
	detectionmgr.RegisterDefaultHandlers(detector)
	detectionmgr.InitCustomDetectionsManager(detector)
	go detectionmgr.StreamLogs(detector)
	if StillInSetup() {
		return
	}
	logger.Detection.Info("Detector loaded successfully")
}

func InitCodeServer() {
	codeserver.InitCodeServer()
}

// InitBundler initialized the onboard bundled assets for the web UI
func InitVirtFS(v1uiFS embed.FS, v2uiFS embed.FS, twoboxFS embed.FS) {
	config.SetV1UIFS(v1uiFS)
	config.SetV2UIFS(v2uiFS)
	config.SetTWOBOXFS(twoboxFS)
}

// InitBackupMgr initializes the backup manager
func InitBackupMgr() {
	backupmgr.InitBackupMgr()
}

func PrintConfigDetails() {
	if StillInSetup() {
		return
	}
	logger.Config.Debug("Gameserver config values loaded")
	logger.Config.Debug("---- GENERAL CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("GameBranch: %s", config.GetGameBranch()))
	logger.Config.Debug("IsDiscordEnabled: " + strconv.FormatBool(config.GetIsDiscordEnabled()))
	logger.Config.Debug("IsCleanupEnabled: " + strconv.FormatBool(config.GetIsCleanupEnabled()))
	logger.Config.Debug("IsDebugMode (pprof Server): " + strconv.FormatBool(config.GetIsDebugMode()))
	logger.Config.Debug("IsFirstTimeSetup: " + strconv.FormatBool(config.GetIsFirstTimeSetup()))

	logger.Config.Debug("---- DISCORD CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BlackListFilePath: %s", config.GetBlackListFilePath()))
	logger.Config.Debug(fmt.Sprintf("ConnectionListChannelID: %s", config.GetConnectionListChannelID()))
	logger.Config.Debug(fmt.Sprintf("ControlChannelID: %s", config.GetControlChannelID()))
	logger.Config.Debug(fmt.Sprintf("ControlPanelChannelID: %s", config.GetControlPanelChannelID()))
	logger.Config.Debug(fmt.Sprintf("DiscordCharBufferSize: %d", config.GetDiscordCharBufferSize()))
	logger.Config.Debug(fmt.Sprintf("DiscordToken: %s", config.GetDiscordToken()))
	logger.Config.Debug(fmt.Sprintf("ErrorChannelID: %s", config.GetErrorChannelID()))
	logger.Config.Debug(fmt.Sprintf("IsDiscordEnabled: %v", config.GetIsDiscordEnabled()))
	logger.Config.Debug(fmt.Sprintf("LogChannelID: %s", config.GetLogChannelID()))
	logger.Config.Debug(fmt.Sprintf("LogMessageBuffer: %s", config.GetLogMessageBuffer()))
	logger.Config.Debug(fmt.Sprintf("SaveChannelID: %s", config.GetSaveChannelID()))
	logger.Config.Debug(fmt.Sprintf("StatusChannelID: %s", config.GetStatusChannelID()))

	logger.Config.Debug("---- BACKUP CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BackupKeepLastN: %d", config.GetBackupKeepLastN()))
	logger.Config.Debug(fmt.Sprintf("BackupKeepDailyFor: %s", config.GetBackupKeepDailyFor()))
	logger.Config.Debug(fmt.Sprintf("BackupKeepWeeklyFor: %s", config.GetBackupKeepWeeklyFor()))
	logger.Config.Debug(fmt.Sprintf("BackupKeepMonthlyFor: %s", config.GetBackupKeepMonthlyFor()))
	logger.Config.Debug(fmt.Sprintf("BackupCleanupInterval: %s", config.GetBackupCleanupInterval()))
	logger.Config.Debug(fmt.Sprintf("ConfiguredBackupDir: %s", config.GetConfiguredBackupDir()))
	logger.Config.Debug(fmt.Sprintf("ConfiguredSafeBackupDir: %s", config.GetConfiguredSafeBackupDir()))
	logger.Config.Debug(fmt.Sprintf("BackupWaitTime: %s", config.GetBackupWaitTime()))

	logger.Config.Debug("---- AUTHENTICATION CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AuthTokenLifetime: %d", config.GetAuthTokenLifetime()))
	logger.Config.Debug(fmt.Sprintf("JwtKey: %s", config.GetJwtKey()))

	logger.Config.Debug("---- SSUI MISC VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("Version: %s", config.Version))

	logger.Config.Debug("----  UPDATER CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AllowPrereleaseUpdates: %v", config.GetAllowPrereleaseUpdates()))
	logger.Config.Debug(fmt.Sprintf("AllowMajorUpdates: %v", config.GetAllowMajorUpdates()))
	logger.Config.Debug(fmt.Sprintf("IsUpdateEnabled: %v", config.GetIsUpdateEnabled()))

	logger.Config.Debug("----  SSCM CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("SSCMFilePath: %s", config.GetSSCMFilePath()))
	logger.Config.Debug(fmt.Sprintf("IsSSCMEnabled: %v", config.GetIsSSCMEnabled()))
}

// StillInSetup returns the current state ("Loading" or "Reloading") and optionally a setup completion status.
func StillInSetup() bool {
	setup.V6setupMutex.Lock()
	defer setup.V6setupMutex.Unlock()
	return !setup.IsSetupComplete
}

func RestartBackend() {
	setup.RestartMySelf()
}
