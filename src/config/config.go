package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	// All configuration variables can be found in vars.go
	Version = "6.4.6"
	Branch  = "v6"
)

type DiscordConfig struct {
	DiscordToken            string `json:"DiscordToken"`
	ControlChannelID        string `json:"ControlChannelID"`
	StatusChannelID         string `json:"StatusChannelID"`
	ConnectionListChannelID string `json:"ConnectionListChannelID"`
	LogChannelID            string `json:"LogChannelID"`
	SaveChannelID           string `json:"SaveChannelID"`
	ControlPanelChannelID   string `json:"ControlPanelChannelID"`
	ErrorChannelID          string `json:"ErrorChannelID"`
	DiscordCharBufferSize   int    `json:"DiscordCharBufferSize"`
	IsDiscordEnabled        *bool  `json:"IsDiscordEnabled"`
}

type AuthConfig struct {
	Users             map[string]string `json:"Users"`
	UserLevels        map[string]string `json:"UserLevels"`
	AuthEnabled       *bool             `json:"AuthEnabled"`
	JwtKey            string            `json:"JwtKey"`
	AuthTokenLifetime int               `json:"AuthTokenLifetime"`
}

type ServerConfig struct {
	RunfileGame         string `json:"RunfileGame"` // Remove this once there is a better way to handle this
	BackendEndpointIP   string `json:"BackendEndpointIP"`
	BackendEndpointPort string `json:"BackendEndpointPort"`
	GameBranch          string `json:"GameBranch"`
}

type LoggingConfig struct {
	Debug             *bool    `json:"Debug"` //pprof
	CreateSSUILogFile *bool    `json:"CreateSSUILogFile"`
	LogLevel          int      `json:"LogLevel"`
	LegacyLogFile     string   `json:"LegacyLogFile"`
	SubsystemFilters  []string `json:"SubsystemFilters"`
}

type UpdateConfig struct {
	IsUpdateEnabled        *bool `json:"IsUpdateEnabled"`
	AllowPrereleaseUpdates *bool `json:"AllowPrereleaseUpdates"`
	AllowMajorUpdates      *bool `json:"AllowMajorUpdates"`
}

type BackupConfig struct {
	BackupContentDir     string        `json:"BackupContentDir"`
	BackupsStoreDir      string        `json:"BackupsStoreDir"`
	BackupLoopInterval   time.Duration `json:"BackupLoopInterval"`
	BackupMode           string        `json:"BackupMode"`
	BackupMaxFileSize    int64         `json:"BackupMaxFileSize"`
	BackupUseCompression *bool         `json:"BackupUseCompression"`
	BackupKeepSnapshot   *bool         `json:"BackupKeepSnapshot"`
}

type SystemConfig struct {
	BlackListFilePath   string `json:"BlackListFilePath"`
	IsSSCMEnabled       *bool  `json:"IsSSCMEnabled"`
	IsFirstTimeSetup    *bool  `json:"IsFirstTimeSetup"`
	IsCodeServerEnabled *bool  `json:"IsCodeServerEnabled"`
	IsConsoleEnabled    *bool  `json:"IsConsoleEnabled"`
}

type JsonConfig struct {
	// A T T E N T I O N : SaveConfig in func applyConfig (below) M U S T reflect this struct!

	Discord DiscordConfig `json:"Discord"`
	Auth    AuthConfig    `json:"Auth"`
	Server  ServerConfig  `json:"Server"`
	Logging LoggingConfig `json:"Logging"`
	Updates UpdateConfig  `json:"Updates"`
	Backup  BackupConfig  `json:"Backup"`
	System  SystemConfig  `json:"System"`
}

// LoadConfig loads and initializes the configuration
func LoadConfig() (*JsonConfig, error) {
	ConfigMu.Lock()

	var jsonConfig JsonConfig
	file, err := os.Open(ConfigPath)
	if err == nil {
		// File exists, proceed to decode it
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&jsonConfig); err != nil {
			return nil, fmt.Errorf("failed to decode config: %v", err)
		}
	} else if os.IsNotExist(err) {
		// File is missing, log it and proceed with defaults (probably first time setup)
		fmt.Println("config file was not found, proceeding with defaults.")
	} else {
		// Other errors (e.g., permissions), fail immediately
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	ConfigMu.Unlock()
	// Apply configuration
	applyConfig(&jsonConfig)

	return &jsonConfig, nil
}

// Environment variables are included solely for automated testing to enable configuration injection during CI/CD pipelines or manual testing.
// They are NOT intended for production use and are undocumented in v6 to discourage reliance.
// The implementation is confusing because:
// 1. Env vars are only read as a fallback during initial configuration and then written to the JSON config file.
// 2. Once written to JSON, env vars are ignored on subsequent runs, leading to confusing behavior.
// Use JSON configuration for reliable and persistent settings.
// applyConfig applies the configuration with JSON -> env -> fallback hierarchy.
func applyConfig(cfg *JsonConfig) {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	// Server Config
	GameBranch = getString(cfg.Server.GameBranch, "GAME_BRANCH", "public")
	BackendEndpointPort = getString(cfg.Server.BackendEndpointPort, "BACKEND_ENDPOINT_PORT", "8443")
	BackendEndpointIP = getString(cfg.Server.BackendEndpointIP, "BACKEND_ENDPOINT_IP", "0.0.0.0")
	if RunfileGame == "" {
		RunfileGame = getString(cfg.Server.RunfileGame, "RUNFILE_GAME", "")
	}

	// Discord Config
	DiscordToken = getString(cfg.Discord.DiscordToken, "DISCORD_TOKEN", "")
	ControlChannelID = getString(cfg.Discord.ControlChannelID, "CONTROL_CHANNEL_ID", "")
	StatusChannelID = getString(cfg.Discord.StatusChannelID, "STATUS_CHANNEL_ID", "")
	ConnectionListChannelID = getString(cfg.Discord.ConnectionListChannelID, "CONNECTION_LIST_CHANNEL_ID", "")
	LogChannelID = getString(cfg.Discord.LogChannelID, "LOG_CHANNEL_ID", "")
	SaveChannelID = getString(cfg.Discord.SaveChannelID, "SAVE_CHANNEL_ID", "")
	ControlPanelChannelID = getString(cfg.Discord.ControlPanelChannelID, "CONTROL_PANEL_CHANNEL_ID", "")
	ErrorChannelID = getString(cfg.Discord.ErrorChannelID, "ERROR_CHANNEL_ID", "")
	DiscordCharBufferSize = getInt(cfg.Discord.DiscordCharBufferSize, "DISCORD_CHAR_BUFFER_SIZE", 1000)

	isDiscordEnabledVal := getBool(cfg.Discord.IsDiscordEnabled, "IS_DISCORD_ENABLED", false)
	IsDiscordEnabled = isDiscordEnabledVal
	cfg.Discord.IsDiscordEnabled = &isDiscordEnabledVal

	// Auth Config
	Users = getMap(cfg.Auth.Users, "SSUI_USERS", map[string]string{})
	UserLevels = getMap(cfg.Auth.UserLevels, "SSUI_USER_LEVELS", map[string]string{})

	authEnabledVal := getBool(cfg.Auth.AuthEnabled, "SSUI_AUTH_ENABLED", false)
	AuthEnabled = authEnabledVal
	cfg.Auth.AuthEnabled = &authEnabledVal

	JwtKey = getString(cfg.Auth.JwtKey, "SSUI_JWT_KEY", generateJwtKey())
	AuthTokenLifetime = getInt(cfg.Auth.AuthTokenLifetime, "SSUI_AUTH_TOKEN_LIFETIME", 1440)

	// Logging Config
	debugVal := getBool(cfg.Logging.Debug, "DEBUG", false)
	IsDebugMode = debugVal
	cfg.Logging.Debug = &debugVal

	createSSUILogFileVal := getBool(cfg.Logging.CreateSSUILogFile, "CREATE_SSUI_LOGFILE", false)
	CreateSSUILogFile = createSSUILogFileVal
	cfg.Logging.CreateSSUILogFile = &createSSUILogFileVal

	LogLevel = getInt(cfg.Logging.LogLevel, "LOG_LEVEL", 20)
	LegacyLogFile = getString(cfg.Logging.LegacyLogFile, "LEGACY_LOG_FILE", "")
	SubsystemFilters = getStringSlice(cfg.Logging.SubsystemFilters, "SUBSYSTEM_FILTERS", []string{})

	// Update Config
	isUpdateEnabledVal := getBool(cfg.Updates.IsUpdateEnabled, "IS_UPDATE_ENABLED", true)
	IsUpdateEnabled = isUpdateEnabledVal
	cfg.Updates.IsUpdateEnabled = &isUpdateEnabledVal

	allowPrereleaseUpdatesVal := getBool(cfg.Updates.AllowPrereleaseUpdates, "ALLOW_PRERELEASE_UPDATES", false)
	AllowPrereleaseUpdates = allowPrereleaseUpdatesVal
	cfg.Updates.AllowPrereleaseUpdates = &allowPrereleaseUpdatesVal

	allowMajorUpdatesVal := getBool(cfg.Updates.AllowMajorUpdates, "ALLOW_MAJOR_UPDATES", false)
	AllowMajorUpdates = allowMajorUpdatesVal
	cfg.Updates.AllowMajorUpdates = &allowMajorUpdatesVal

	// System Config
	BlackListFilePath = getString(cfg.System.BlackListFilePath, "BLACKLIST_FILE_PATH", "./Blacklist.txt")

	isSSCMEnabledVal := getBool(cfg.System.IsSSCMEnabled, "IS_SSCM_ENABLED", false)
	IsSSCMEnabled = isSSCMEnabledVal
	cfg.System.IsSSCMEnabled = &isSSCMEnabledVal

	IsFirstTimeSetup = getBool(cfg.System.IsFirstTimeSetup, "IS_FIRST_TIME_SETUP", true)
	cfg.System.IsFirstTimeSetup = &IsFirstTimeSetup

	IsCodeServerEnabled = getBool(cfg.System.IsCodeServerEnabled, "IS_CODE_SERVER_ENABLED", false)
	cfg.System.IsCodeServerEnabled = &IsCodeServerEnabled

	IsConsoleEnabled = getBool(cfg.System.IsConsoleEnabled, "IS_CONSOLE_ENABLED", false)
	cfg.System.IsConsoleEnabled = &IsConsoleEnabled

	// Backup Manager v3 Settings
	BackupContentDir = getString(cfg.Backup.BackupContentDir, "BACKUP_CONTENT_DIR", UIModFolder+"backups/content")
	BackupsStoreDir = getString(cfg.Backup.BackupsStoreDir, "STORED_BACKUPS_DIR", UIModFolder+"backups/storedBackups")
	BackupLoopInterval = getDuration(cfg.Backup.BackupLoopInterval, "BACKUP_LOOP_INTERVAL", time.Hour)
	BackupMode = getString(cfg.Backup.BackupMode, "BACKUP_MODE", "tar")
	BackupMaxFileSize = getInt64(cfg.Backup.BackupMaxFileSize, "MAX_FILE_SIZE", 20*1024*1024*1024)
	BackupUseCompression = getBool(cfg.Backup.BackupUseCompression, "USE_COMPRESSION", true)
	BackupKeepSnapshot = getBool(cfg.Backup.BackupKeepSnapshot, "KEEP_SNAPSHOT", false)
}

// saveConfig M U S T be called while holding a lock on ConfigMu! Accepts an optional deferred action to run after successfully saving the config
func saveConfig(deferredAction ...DeferredAction) error {

	cfg := JsonConfig{
		Server: ServerConfig{
			RunfileGame:         RunfileGame,
			BackendEndpointIP:   BackendEndpointIP,
			BackendEndpointPort: BackendEndpointPort,
			GameBranch:          GameBranch,
		},
		Discord: DiscordConfig{
			DiscordToken:            DiscordToken,
			ControlChannelID:        ControlChannelID,
			StatusChannelID:         StatusChannelID,
			ConnectionListChannelID: ConnectionListChannelID,
			LogChannelID:            LogChannelID,
			SaveChannelID:           SaveChannelID,
			ControlPanelChannelID:   ControlPanelChannelID,
			ErrorChannelID:          ErrorChannelID,
			DiscordCharBufferSize:   DiscordCharBufferSize,
			IsDiscordEnabled:        &IsDiscordEnabled,
		},
		Auth: AuthConfig{
			Users:             Users,
			UserLevels:        UserLevels,
			AuthEnabled:       &AuthEnabled,
			JwtKey:            JwtKey,
			AuthTokenLifetime: AuthTokenLifetime,
		},
		Logging: LoggingConfig{
			Debug:             &IsDebugMode,
			CreateSSUILogFile: &CreateSSUILogFile,
			LogLevel:          LogLevel,
			LegacyLogFile:     LegacyLogFile,
			SubsystemFilters:  SubsystemFilters,
		},
		Updates: UpdateConfig{
			IsUpdateEnabled:        &IsUpdateEnabled,
			AllowPrereleaseUpdates: &AllowPrereleaseUpdates,
			AllowMajorUpdates:      &AllowMajorUpdates,
		},
		Backup: BackupConfig{
			BackupContentDir:     BackupContentDir,
			BackupsStoreDir:      BackupsStoreDir,
			BackupLoopInterval:   BackupLoopInterval,
			BackupMode:           BackupMode,
			BackupMaxFileSize:    BackupMaxFileSize,
			BackupUseCompression: &BackupUseCompression,
			BackupKeepSnapshot:   &BackupKeepSnapshot,
		},
		System: SystemConfig{
			BlackListFilePath:   BlackListFilePath,
			IsSSCMEnabled:       &IsSSCMEnabled,
			IsFirstTimeSetup:    &IsFirstTimeSetup,
			IsCodeServerEnabled: &IsCodeServerEnabled,
			IsConsoleEnabled:    &IsConsoleEnabled,
		},
	}

	file, err := os.Create(ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&cfg); err != nil {
		return fmt.Errorf("failed to encode config: %v", err)
	}

	// Run deferred action after saving, if provided
	if len(deferredAction) > 0 {
		runDeferredAction(deferredAction[0])
	}

	return nil
}
