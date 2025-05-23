package config

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

// GetWorldName returns the WorldName
func GetWorldName() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return WorldName
}

// GetBackupWorldName returns the BackupWorldName
func GetBackupWorldName() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupWorldName
}

// GetIsDebugMode returns the IsDebugMode
func GetIsDebugMode() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsDebugMode
}

// GetCreateSSUILogFile returns the CreateSSUILogFile
func GetCreateSSUILogFile() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return CreateSSUILogFile
}

// GetLogLevel returns the LogLevel
func GetLogLevel() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogLevel
}

// GetLogMessageBuffer returns the LogMessageBuffer
func GetLogMessageBuffer() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogMessageBuffer
}

// GetIsFirstTimeSetup returns the IsFirstTimeSetup
func GetIsFirstTimeSetup() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsFirstTimeSetup
}

// GetBufferFlushTicker returns the BufferFlushTicker
func GetBufferFlushTicker() *time.Ticker {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BufferFlushTicker
}

// GetSSEMessageBufferSize returns the SSEMessageBufferSize
func GetSSEMessageBufferSize() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSEMessageBufferSize
}

// GetMaxSSEConnections returns the MaxSSEConnections
func GetMaxSSEConnections() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return MaxSSEConnections
}

// GetGameServerAppID returns the GameServerAppID
func GetGameServerAppID() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return GameServerAppID
}

// GetGameBranch returns the GameBranch
func GetGameBranch() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return GameBranch
}

// GetSubsystemFilters returns the SubsystemFilters
func GetSubsystemFilters() []string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SubsystemFilters
}

// GetGameServerUUID returns the GameServerUUID
func GetGameServerUUID() uuid.UUID {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return GameServerUUID
}

// GetBackendEndpointPort returns the BackendEndpointPort
func GetBackendEndpointPort() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackendEndpointPort
}

// GetBackendEndpointIP returns the BackendEndpointIP
func GetBackendEndpointIP() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackendEndpointIP
}

// GetDiscordToken returns the DiscordToken
func GetDiscordToken() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return DiscordToken
}

// GetDiscordSession returns the DiscordSession
func GetDiscordSession() *discordgo.Session {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return DiscordSession
}

// GetIsDiscordEnabled returns the IsDiscordEnabled
func GetIsDiscordEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsDiscordEnabled
}

// GetControlChannelID returns the ControlChannelID
func GetControlChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ControlChannelID
}

// GetStatusChannelID returns the StatusChannelID
func GetStatusChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return StatusChannelID
}

// GetLogChannelID returns the LogChannelID
func GetLogChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogChannelID
}

// GetErrorChannelID returns the ErrorChannelID
func GetErrorChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ErrorChannelID
}

// GetConnectionListChannelID returns the ConnectionListChannelID
func GetConnectionListChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConnectionListChannelID
}

// GetSaveChannelID returns the SaveChannelID
func GetSaveChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SaveChannelID
}

// GetControlPanelChannelID returns the ControlPanelChannelID
func GetControlPanelChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ControlPanelChannelID
}

// GetDiscordCharBufferSize returns the DiscordCharBufferSize
func GetDiscordCharBufferSize() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return DiscordCharBufferSize
}

// GetControlMessageID returns the ControlMessageID
func GetControlMessageID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ControlMessageID
}

// GetExceptionMessageID returns the ExceptionMessageID
func GetExceptionMessageID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ExceptionMessageID
}

// GetBlackListFilePath returns the BlackListFilePath
func GetBlackListFilePath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BlackListFilePath
}

// GetIsCleanupEnabled returns the IsCleanupEnabled
func GetIsCleanupEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsCleanupEnabled
}

// GetBackupKeepLastN returns the BackupKeepLastN
func GetBackupKeepLastN() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepLastN
}

// GetBackupKeepDailyFor returns the BackupKeepDailyFor
func GetBackupKeepDailyFor() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepDailyFor
}

// GetBackupKeepWeeklyFor returns the BackupKeepWeeklyFor
func GetBackupKeepWeeklyFor() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepWeeklyFor
}

// GetBackupKeepMonthlyFor returns the BackupKeepMonthlyFor
func GetBackupKeepMonthlyFor() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepMonthlyFor
}

// GetBackupCleanupInterval returns the BackupCleanupInterval
func GetBackupCleanupInterval() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupCleanupInterval
}

// GetConfiguredBackupDir returns the ConfiguredBackupDir
func GetConfiguredBackupDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConfiguredBackupDir
}

// GetConfiguredSafeBackupDir returns the ConfiguredSafeBackupDir
func GetConfiguredSafeBackupDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConfiguredSafeBackupDir
}

// GetBackupWaitTime returns the BackupWaitTime
func GetBackupWaitTime() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupWaitTime
}

// GetAuthEnabled returns the AuthEnabled
func GetAuthEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AuthEnabled
}

// GetJwtKey returns the JwtKey
func GetJwtKey() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return JwtKey
}

// GetAuthTokenLifetime returns the AuthTokenLifetime
func GetAuthTokenLifetime() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AuthTokenLifetime
}

// GetUsers returns the Users
func GetUsers() map[string]string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return Users
}

// GetIsUpdateEnabled returns the IsUpdateEnabled
func GetIsUpdateEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsUpdateEnabled
}

// GetAllowPrereleaseUpdates returns the AllowPrereleaseUpdates
func GetAllowPrereleaseUpdates() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AllowPrereleaseUpdates
}

// GetAllowMajorUpdates returns the AllowMajorUpdates
func GetAllowMajorUpdates() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AllowMajorUpdates
}

// GetIsSSCMEnabled returns the IsSSCMEnabled
func GetIsSSCMEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsSSCMEnabled
}

// GetRunfileGame returns the RunfileGame
func GetRunfileGame() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunfileGame
}

// GetSteamCMDLinuxDir returns the SteamCMDLinuxDir
func GetSteamCMDLinuxDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SteamCMDLinuxDir
}

// GetSteamCMDWindowsDir returns the SteamCMDWindowsDir
func GetSteamCMDWindowsDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SteamCMDWindowsDir
}

// GetSteamCMDLinuxURL returns the SteamCMDLinuxURL
func GetSteamCMDLinuxURL() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SteamCMDLinuxURL
}

// GetSteamCMDWindowsURL returns the SteamCMDWindowsURL
func GetSteamCMDWindowsURL() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SteamCMDWindowsURL
}

// GetTLSCertPath returns the TLSCertPath
func GetTLSCertPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return TLSCertPath
}

// GetTLSKeyPath returns the TLSKeyPath
func GetTLSKeyPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return TLSKeyPath
}

// GetConfigPath returns the ConfigPath
func GetConfigPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConfigPath
}

// GetCustomDetectionsFilePath returns the CustomDetectionsFilePath
func GetCustomDetectionsFilePath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return CustomDetectionsFilePath
}

// GetLogFolder returns the LogFolder
func GetLogFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogFolder
}

// GetUIModFolder returns the UIModFolder
func GetUIModFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return UIModFolder
}

// GetTwoBoxFormFolder returns the TwoBoxFormFolder
func GetTwoBoxFormFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return TwoBoxFormFolder
}

// GetConfigHtmlPath returns the ConfigHtmlPath
func GetConfigHtmlPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConfigHtmlPath
}

// GetDetectionManagerHtmlPath returns the DetectionManagerHtmlPath
func GetDetectionManagerHtmlPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return DetectionManagerHtmlPath
}

// GetTwoBoxFormHtmlPath returns the TwoBoxFormHtmlPath
func GetTwoBoxFormHtmlPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return TwoBoxFormHtmlPath
}

// GetIndexHtmlPath returns the IndexHtmlPath
func GetIndexHtmlPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IndexHtmlPath
}

// GetSSCMWebDir returns the SSCMWebDir
func GetSSCMWebDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSCMWebDir
}

// GetSSCMFilePath returns the SSCMFilePath
func GetSSCMFilePath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSCMFilePath
}

// GetSSCMPluginDir returns the SSCMPluginDir
func GetSSCMPluginDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSCMPluginDir
}

// GetRunFilesFolder returns the RunFilesFolder
func GetRunFilesFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunFilesFolder
}

func GetLegacyLogFile() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LegacyLogFile
}

func GetRunfilesFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunFilesFolder
}

func GetBackendVersion() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return Version
}
