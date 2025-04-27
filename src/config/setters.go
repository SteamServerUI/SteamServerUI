package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SetWorldName sets the WorldName with validation
func SetWorldName(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("world name cannot be empty")
	}

	WorldName = value
	return SaveConfig()
}

// SetBackupWorldName sets the BackupWorldName with validation
func SetBackupWorldName(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("backup world name cannot be empty")
	}

	BackupWorldName = value
	return SaveConfig()
}

// SetIsDebugMode sets the IsDebugMode with validation
func SetIsDebugMode(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDebugMode = value
	return SaveConfig()
}

// SetCreateSSUILogFile sets the CreateSSUILogFile with validation
func SetCreateSSUILogFile(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	CreateSSUILogFile = value
	return SaveConfig()
}

// SetLogLevel sets the LogLevel with validation
func SetLogLevel(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("log level cannot be negative")
	}

	LogLevel = value
	return SaveConfig()
}

// SetLogMessageBuffer sets the LogMessageBuffer with validation
func SetLogMessageBuffer(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogMessageBuffer = value
	return SaveConfig()
}

// SetIsFirstTimeSetup sets the IsFirstTimeSetup with validation
func SetIsFirstTimeSetup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsFirstTimeSetup = value
	return SaveConfig()
}

// SetBufferFlushTicker sets the BufferFlushTicker
func SetBufferFlushTicker(value *time.Ticker) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BufferFlushTicker = value
	return SaveConfig()
}

// SetSSEMessageBufferSize sets the SSEMessageBufferSize with validation
func SetSSEMessageBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("SSE message buffer size must be positive")
	}

	SSEMessageBufferSize = value
	return SaveConfig()
}

// SetMaxSSEConnections sets the MaxSSEConnections with validation
func SetMaxSSEConnections(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("max SSE connections must be positive")
	}

	MaxSSEConnections = value
	return SaveConfig()
}

// SetGameServerAppID sets the GameServerAppID with validation
func SetGameServerAppID(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("game server app ID must be positive")
	}

	GameServerAppID = value
	return SaveConfig()
}

// SetGameBranch sets the GameBranch with validation
func SetGameBranch(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("game branch cannot be empty")
	}

	GameBranch = value
	return SaveConfig()
}

// SetSubsystemFilters sets the SubsystemFilters with validation
func SetSubsystemFilters(value []string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	for _, v := range value {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("subsystem filter cannot be empty")
		}
	}

	SubsystemFilters = value
	return SaveConfig()
}

// SetGameServerUUID sets the GameServerUUID with validation
func SetGameServerUUID(value uuid.UUID) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value == uuid.Nil {
		return fmt.Errorf("game server UUID cannot be nil")
	}

	GameServerUUID = value
	return SaveConfig()
}

// SetBackendEndpointPort sets the BackendEndpointPort with validation
func SetBackendEndpointPort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("backend endpoint port cannot be empty")
	}

	BackendEndpointPort = value
	return SaveConfig()
}

// SetBackendEndpointIP sets the BackendEndpointIP with validation
func SetBackendEndpointIP(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("backend endpoint IP cannot be empty")
	}

	BackendEndpointIP = value
	return SaveConfig()
}

// SetDiscordToken sets the DiscordToken with validation
func SetDiscordToken(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	DiscordToken = value
	return SaveConfig()
}

// SetIsDiscordEnabled sets the IsDiscordEnabled with validation
func SetIsDiscordEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDiscordEnabled = value
	return SaveConfig()
}

// SetControlChannelID sets the ControlChannelID with validation
func SetControlChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlChannelID = value
	return SaveConfig()
}

// SetStatusChannelID sets the StatusChannelID with validation
func SetStatusChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StatusChannelID = value
	return SaveConfig()
}

// SetLogChannelID sets the LogChannelID with validation
func SetLogChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogChannelID = value
	return SaveConfig()
}

// SetErrorChannelID sets the ErrorChannelID with validation
func SetErrorChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ErrorChannelID = value
	return SaveConfig()
}

// SetConnectionListChannelID sets the ConnectionListChannelID with validation
func SetConnectionListChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ConnectionListChannelID = value
	return SaveConfig()
}

// SetSaveChannelID sets the SaveChannelID with validation
func SetSaveChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveChannelID = value
	return SaveConfig()
}

// SetControlPanelChannelID sets the ControlPanelChannelID with validation
func SetControlPanelChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlPanelChannelID = value
	return SaveConfig()
}

// SetDiscordCharBufferSize sets the DiscordCharBufferSize with validation
func SetDiscordCharBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("discord char buffer size must be positive")
	}

	DiscordCharBufferSize = value
	return SaveConfig()
}

// SetControlMessageID sets the ControlMessageID with validation
func SetControlMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlMessageID = value
	return SaveConfig()
}

// SetExceptionMessageID sets the ExceptionMessageID with validation
func SetExceptionMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExceptionMessageID = value
	return SaveConfig()
}

// SetBlackListFilePath sets the BlackListFilePath with validation
func SetBlackListFilePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("blacklist file path cannot be empty")
	}

	BlackListFilePath = value
	return SaveConfig()
}

// SetIsCleanupEnabled sets the IsCleanupEnabled with validation
func SetIsCleanupEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsCleanupEnabled = value
	return SaveConfig()
}

// SetBackupKeepLastN sets the BackupKeepLastN with validation
func SetBackupKeepLastN(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("backup keep last N must be positive")
	}

	BackupKeepLastN = value
	return SaveConfig()
}

// SetBackupKeepDailyFor sets the BackupKeepDailyFor with validation
func SetBackupKeepDailyFor(value time.Duration) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep daily for cannot be negative")
	}

	BackupKeepDailyFor = value
	return SaveConfig()
}

// SetBackupKeepWeeklyFor sets the BackupKeepWeeklyFor with validation
func SetBackupKeepWeeklyFor(value time.Duration) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep weekly for cannot be negative")
	}

	BackupKeepWeeklyFor = value
	return SaveConfig()
}

// SetBackupKeepMonthlyFor sets the BackupKeepMonthlyFor with validation
func SetBackupKeepMonthlyFor(value time.Duration) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep monthly for cannot be negative")
	}

	BackupKeepMonthlyFor = value
	return SaveConfig()
}

// SetBackupCleanupInterval sets the BackupCleanupInterval with validation
func SetBackupCleanupInterval(value time.Duration) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("backup cleanup interval must be positive")
	}

	BackupCleanupInterval = value
	return SaveConfig()
}

// SetConfiguredBackupDir sets the ConfiguredBackupDir with validation
func SetConfiguredBackupDir(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("configured backup dir cannot be empty")
	}

	ConfiguredBackupDir = value
	return SaveConfig()
}

// SetConfiguredSafeBackupDir sets the ConfiguredSafeBackupDir with validation
func SetConfiguredSafeBackupDir(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("configured safe backup dir cannot be empty")
	}

	ConfiguredSafeBackupDir = value
	return SaveConfig()
}

// SetBackupWaitTime sets the BackupWaitTime with validation
func SetBackupWaitTime(value time.Duration) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup wait time cannot be negative")
	}

	BackupWaitTime = value
	return SaveConfig()
}

// SetAuthEnabled sets the AuthEnabled with validation
func SetAuthEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AuthEnabled = value
	return SaveConfig()
}

// SetJwtKey sets the JwtKey with validation
func SetJwtKey(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if len(value) < 32 {
		return fmt.Errorf("JWT key must be at least 32 bytes")
	}

	JwtKey = value
	return SaveConfig()
}

// SetAuthTokenLifetime sets the AuthTokenLifetime with validation
func SetAuthTokenLifetime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("auth token lifetime must be positive")
	}

	AuthTokenLifetime = value
	return SaveConfig()
}

// SetUsers sets the Users with validation
func SetUsers(value map[string]string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	for k, v := range value {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			return fmt.Errorf("user key or value cannot be empty")
		}
	}

	Users = value
	return SaveConfig()
}

// SetIsUpdateEnabled sets the IsUpdateEnabled with validation
func SetIsUpdateEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsUpdateEnabled = value
	return SaveConfig()
}

// SetAllowPrereleaseUpdates sets the AllowPrereleaseUpdates with validation
func SetAllowPrereleaseUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowPrereleaseUpdates = value
	return SaveConfig()
}

// SetAllowMajorUpdates sets the AllowMajorUpdates with validation
func SetAllowMajorUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowMajorUpdates = value
	return SaveConfig()
}

// SetIsSSCMEnabled sets the IsSSCMEnabled with validation
func SetIsSSCMEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsSSCMEnabled = value
	return SaveConfig()
}

// SetRunfileGame sets the RunfileGame with validation
func SetRunfileGame(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("runfile game cannot be empty")
	}

	RunfileGame = value
	return SaveConfig() // Deferred action below is not implemented yet
	//return SaveConfig(func() {
	//	loader.ReloadSomethingRunfileRelated() // Deferred action to reload Runfile-related things, not implemented yet because loader imports config already.
	//  I am unsure how to resolve this import cycle. Maybe emit a signal instead?
	//})
}
