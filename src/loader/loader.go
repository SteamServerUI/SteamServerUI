// loader.go
package loader

import (
	"fmt"
	"strconv"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/argmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/detectionmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/discordbot"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
)

func ReloadAll() {
	ReloadConfig()
	ReloadBackupManager()
	ReloadDiscordBot()
	ReloadRunfile()
}

func ReloadConfig() {
	if _, err := config.LoadConfig(); err != nil {
		logger.Core.Error("Failed to load config: " + err.Error())
		return
	}

	logger.Core.Info("Config reloaded successfully")

	if config.IsSSCMEnabled {
		setup.InstallSSCM()
	}

	PrintConfigDetails()
}

func ReloadBackupManager() {
	logger.Backup.Info("Backup manager cannot be initialized in SteamServerUI")
	//if err := backupmgr.ReloadBackupManagerFromConfig(); err != nil {
	//	logger.Backup.Error("Failed to reload backup manager: " + err.Error())
	//	return
	//}
	//logger.Backup.Info("Backup manager reloaded successfully")
}

func ReloadDiscordBot() {
	if config.IsDiscordEnabled {
		go discordbot.InitializeDiscordBot()
		logger.Discord.Info("Discord bot reloaded successfully")
	}
}

func ReloadRunfile() {
	if err := argmgr.LoadRunfile(config.RunfileGame, config.RunFilesFolder); err != nil {
		logger.Runfile.Error("Failed to reload runfile: " + err.Error())
		return
	}
	logger.Runfile.Info("Runfile reloaded successfully")
}

// The detector should NOT be reloaded, as it is a singleton. Instead, dynamic changes come in via the custom detections manager.
func InitDetector() {
	detector := detectionmgr.Start()
	detectionmgr.RegisterDefaultHandlers(detector)
	detectionmgr.InitCustomDetectionsManager(detector)
	go detectionmgr.StreamLogs(detector)
	logger.Detection.Info("Detector loaded successfully")
}

func PrintConfigDetails() {
	logger.Config.Debug("Gameserver config values loaded")
	logger.Config.Debug("---- GENERAL CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("GameBranch: %s", config.GameBranch))
	logger.Config.Debug("IsDiscordEnabled: " + strconv.FormatBool(config.IsDiscordEnabled))
	logger.Config.Debug("IsCleanupEnabled: " + strconv.FormatBool(config.IsCleanupEnabled))
	logger.Config.Debug("IsDebugMode (pprof Server): " + strconv.FormatBool(config.IsDebugMode))
	logger.Config.Debug("IsFirstTimeSetup: " + strconv.FormatBool(config.IsFirstTimeSetup))

	logger.Config.Debug("---- DISCORD CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BlackListFilePath: %s", config.BlackListFilePath))
	logger.Config.Debug(fmt.Sprintf("ConnectionListChannelID: %s", config.ConnectionListChannelID))
	logger.Config.Debug(fmt.Sprintf("ControlChannelID: %s", config.ControlChannelID))
	logger.Config.Debug(fmt.Sprintf("ControlPanelChannelID: %s", config.ControlPanelChannelID))
	logger.Config.Debug(fmt.Sprintf("DiscordCharBufferSize: %d", config.DiscordCharBufferSize))
	logger.Config.Debug(fmt.Sprintf("DiscordToken: %s", config.DiscordToken))
	logger.Config.Debug(fmt.Sprintf("ErrorChannelID: %s", config.ErrorChannelID))
	logger.Config.Debug(fmt.Sprintf("IsDiscordEnabled: %v", config.IsDiscordEnabled))
	logger.Config.Debug(fmt.Sprintf("LogChannelID: %s", config.LogChannelID))
	logger.Config.Debug(fmt.Sprintf("LogMessageBuffer: %s", config.LogMessageBuffer))
	logger.Config.Debug(fmt.Sprintf("SaveChannelID: %s", config.SaveChannelID))
	logger.Config.Debug(fmt.Sprintf("StatusChannelID: %s", config.StatusChannelID))

	logger.Config.Debug("---- BACKUP CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BackupKeepLastN: %d", config.BackupKeepLastN))
	logger.Config.Debug(fmt.Sprintf("BackupKeepDailyFor: %s", config.BackupKeepDailyFor))
	logger.Config.Debug(fmt.Sprintf("BackupKeepWeeklyFor: %s", config.BackupKeepWeeklyFor))
	logger.Config.Debug(fmt.Sprintf("BackupKeepMonthlyFor: %s", config.BackupKeepMonthlyFor))
	logger.Config.Debug(fmt.Sprintf("BackupCleanupInterval: %s", config.BackupCleanupInterval))
	logger.Config.Debug(fmt.Sprintf("ConfiguredBackupDir: %s", config.ConfiguredBackupDir))
	logger.Config.Debug(fmt.Sprintf("ConfiguredSafeBackupDir: %s", config.ConfiguredSafeBackupDir))
	logger.Config.Debug(fmt.Sprintf("BackupWaitTime: %s", config.BackupWaitTime))

	logger.Config.Debug("---- AUTHENTICATION CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AuthTokenLifetime: %d", config.AuthTokenLifetime))
	logger.Config.Debug(fmt.Sprintf("JwtKey: %s", config.JwtKey))

	logger.Config.Debug("---- SSUI MISC VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("Version: %s", config.Version))

	logger.Config.Debug("----  UPDATER CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AllowPrereleaseUpdates: %v", config.AllowPrereleaseUpdates))
	logger.Config.Debug(fmt.Sprintf("AllowMajorUpdates: %v", config.AllowMajorUpdates))
	logger.Config.Debug(fmt.Sprintf("IsUpdateEnabled: %v", config.IsUpdateEnabled))

	logger.Config.Debug("----  SSCM CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("SSCMFilePath: %s", config.SSCMFilePath))
	logger.Config.Debug(fmt.Sprintf("IsSSCMEnabled: %v", config.IsSSCMEnabled))
}
