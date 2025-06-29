package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Although this is a not a real setter, this function can be used to save the config safely
func SetSaveConfig() error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return saveConfig()
}

// SetIsDebugMode sets the IsDebugMode with validation
func SetIsDebugMode(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDebugMode = value
	return saveConfig()
}

// SetCreateSSUILogFile sets the CreateSSUILogFile with validation
func SetCreateSSUILogFile(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	CreateSSUILogFile = value
	return saveConfig()
}

// SetLogLevel sets the LogLevel with validation
func SetLogLevel(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("log level cannot be negative")
	}

	LogLevel = value
	return saveConfig()
}

// SetLogMessageBuffer sets the LogMessageBuffer with validation
func SetLogMessageBuffer(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogMessageBuffer = value
	return saveConfig()
}

// SetIsFirstTimeSetup sets the IsFirstTimeSetup with validation
func SetIsFirstTimeSetup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsFirstTimeSetup = value
	return saveConfig()
}

// SetBufferFlushTicker sets the BufferFlushTicker
func SetBufferFlushTicker(value *time.Ticker) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BufferFlushTicker = value
	return saveConfig()
}

// SetSSEMessageBufferSize sets the SSEMessageBufferSize with validation
func SetSSEMessageBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("SSE message buffer size must be positive")
	}

	SSEMessageBufferSize = value
	return saveConfig()
}

// SetMaxSSEConnections sets the MaxSSEConnections with validation
func SetMaxSSEConnections(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("max SSE connections must be positive")
	}

	MaxSSEConnections = value
	return saveConfig()
}

// SetGameServerAppID sets the GameServerAppID with validation
func SetGameServerAppID(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("game server app ID must be positive")
	}

	GameServerAppID = value
	return saveConfig()
}

// SetGameBranch sets the GameBranch with validation
func SetGameBranch(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("game branch cannot be empty")
	}

	GameBranch = value
	return saveConfig()
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
	return saveConfig()
}

// SetGameServerUUID sets the GameServerUUID with validation
func SetGameServerUUID(value uuid.UUID) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value == uuid.Nil {
		return fmt.Errorf("game server UUID cannot be nil")
	}

	GameServerUUID = value
	return saveConfig()
}

// SetBackendEndpointPort sets the BackendEndpointPort with validation
func SetBackendEndpointPort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("backend endpoint port cannot be empty")
	}

	BackendEndpointPort = value
	return saveConfig()
}

// SetBackendEndpointIP sets the BackendEndpointIP with validation
func SetBackendEndpointIP(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("backend endpoint IP cannot be empty")
	}

	BackendEndpointIP = value
	return saveConfig()
}

// SetDiscordToken sets the DiscordToken with validation
func SetDiscordToken(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	DiscordToken = value
	return saveConfig()
}

// SetIsDiscordEnabled sets the IsDiscordEnabled with validation
func SetIsDiscordEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDiscordEnabled = value
	return saveConfig()
}

// SetControlChannelID sets the ControlChannelID with validation
func SetControlChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlChannelID = value
	return saveConfig()
}

// SetStatusChannelID sets the StatusChannelID with validation
func SetStatusChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StatusChannelID = value
	return saveConfig()
}

// SetLogChannelID sets the LogChannelID with validation
func SetLogChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogChannelID = value
	return saveConfig()
}

// SetErrorChannelID sets the ErrorChannelID with validation
func SetErrorChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ErrorChannelID = value
	return saveConfig()
}

// SetConnectionListChannelID sets the ConnectionListChannelID with validation
func SetConnectionListChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ConnectionListChannelID = value
	return saveConfig()
}

// SetSaveChannelID sets the SaveChannelID with validation
func SetSaveChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveChannelID = value
	return saveConfig()
}

// SetControlPanelChannelID sets the ControlPanelChannelID with validation
func SetControlPanelChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlPanelChannelID = value
	return saveConfig()
}

// SetDiscordCharBufferSize sets the DiscordCharBufferSize with validation
func SetDiscordCharBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("discord char buffer size must be positive")
	}

	DiscordCharBufferSize = value
	return saveConfig()
}

// SetControlMessageID sets the ControlMessageID with validation
func SetControlMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlMessageID = value
	return saveConfig()
}

// SetExceptionMessageID sets the ExceptionMessageID with validation
func SetExceptionMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExceptionMessageID = value
	return saveConfig()
}

// SetBlackListFilePath sets the BlackListFilePath with validation
func SetBlackListFilePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("blacklist file path cannot be empty")
	}

	BlackListFilePath = value
	return saveConfig()
}

// SetAuthEnabled sets the AuthEnabled with validation
func SetAuthEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AuthEnabled = value
	return saveConfig()
}

// SetJwtKey sets the JwtKey with validation
func SetJwtKey(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if len(value) < 32 {
		return fmt.Errorf("JWT key must be at least 32 bytes")
	}

	JwtKey = value
	return saveConfig()
}

// SetAuthTokenLifetime sets the AuthTokenLifetime with validation
func SetAuthTokenLifetime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("auth token lifetime must be positive")
	}

	AuthTokenLifetime = value
	return saveConfig()
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
	return saveConfig()
}

// SetIsUpdateEnabled sets the IsUpdateEnabled with validation
func SetIsUpdateEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsUpdateEnabled = value
	return saveConfig()
}

// SetAllowPrereleaseUpdates sets the AllowPrereleaseUpdates with validation
func SetAllowPrereleaseUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowPrereleaseUpdates = value
	return saveConfig()
}

// SetAllowMajorUpdates sets the AllowMajorUpdates with validation
func SetAllowMajorUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowMajorUpdates = value
	return saveConfig()
}

// SetIsSSCMEnabled sets the IsSSCMEnabled with validation
func SetIsSSCMEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsSSCMEnabled = value
	return saveConfig()
}

// SetRunfileGame sets the RunfileGame with validation
func SetRunfileGame(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("runfile game cannot be empty")
	}

	RunfileGame = value
	return saveConfig() // Deferred action below is not implemented yet
	//return saveConfig(func() {
	//	loader.ReloadSomethingRunfileRelated() // Deferred action to reload Runfile-related things, not implemented yet because loader imports config already.
	//  I am unsure how to resolve this import cycle. Maybe emit a signal instead?
	//})
}

func SetLegacyLogFile(value string) error {
	if !strings.HasPrefix(value, "./") {
		return fmt.Errorf("legacy log file path must start with './'")
	}

	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LegacyLogFile = value
	return saveConfig()
}

func SetIsCodeServerEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsCodeServerEnabled = value
	return saveConfig()
}

func SetBackupContentDir(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupContentDir = value
	return saveConfig()
}

func SetBackupsStoreDir(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupsStoreDir = value
	return saveConfig()
}

func SetBackupLoopInterval(value time.Duration) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupLoopInterval = value
	return saveConfig()
}

func SetBackupMode(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupMode = value
	return saveConfig()
}

func SetBackupMaxFileSize(value int64) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupMaxFileSize = value
	return saveConfig()
}

func SetBackupUseCompression(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupUseCompression = value
	return saveConfig()
}

func SetBackupKeepSnapshot(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupKeepSnapshot = value
	return saveConfig()
}
