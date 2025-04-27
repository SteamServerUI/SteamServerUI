package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	// All configuration variables can be found in vars.go
	Version              = "6.0.11"
	Branch               = "v6-pre"
	IsSteamServerUIBuild = true
)

type JsonConfig struct {
	RunfileGame             string            `json:"runfileGame"` // Remove this once there is a better way to handle this
	BackendEndpointIP       string            `json:"backendEndpointIP"`
	BackendEndpointPort     string            `json:"backendEndpointPort"`
	DiscordToken            string            `json:"discordToken"`
	ControlChannelID        string            `json:"controlChannelID"`
	StatusChannelID         string            `json:"statusChannelID"`
	ConnectionListChannelID string            `json:"connectionListChannelID"`
	LogChannelID            string            `json:"logChannelID"`
	SaveChannelID           string            `json:"saveChannelID"`
	ControlPanelChannelID   string            `json:"controlPanelChannelID"`
	DiscordCharBufferSize   int               `json:"DiscordCharBufferSize"`
	BlackListFilePath       string            `json:"blackListFilePath"`
	IsDiscordEnabled        *bool             `json:"isDiscordEnabled"`
	ErrorChannelID          string            `json:"errorChannelID"`
	BackupKeepLastN         int               `json:"backupKeepLastN"`
	IsCleanupEnabled        *bool             `json:"isCleanupEnabled"`
	BackupKeepDailyFor      int               `json:"backupKeepDailyFor"`
	BackupKeepWeeklyFor     int               `json:"backupKeepWeeklyFor"`
	BackupKeepMonthlyFor    int               `json:"backupKeepMonthlyFor"`
	BackupCleanupInterval   int               `json:"backupCleanupInterval"`
	BackupWaitTime          int               `json:"backupWaitTime"`
	GameBranch              string            `json:"gameBranch"`
	Users                   map[string]string `json:"users"`       // Map of username to hashed password
	AuthEnabled             *bool             `json:"authEnabled"` // Toggle for enabling/disabling auth
	JwtKey                  string            `json:"JwtKey"`
	AuthTokenLifetime       int               `json:"AuthTokenLifetime"`
	Debug                   *bool             `json:"Debug"` //pprof
	CreateSSUILogFile       *bool             `json:"CreateSSUILogFile"`
	LogLevel                int               `json:"LogLevel"`
	SubsystemFilters        []string          `json:"subsystemFilters"`
	IsUpdateEnabled         *bool             `json:"IsUpdateEnabled"`
	IsSSCMEnabled           *bool             `json:"IsSSCMEnabled"`
	AllowPrereleaseUpdates  *bool             `json:"AllowPrereleaseUpdates"`
	AllowMajorUpdates       *bool             `json:"AllowMajorUpdates"`
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
		// File is missing, log it and proceed with defaults
		fmt.Println("Config file does not exist.")
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

	// SteamServerUI config
	GameBranch = getString(cfg.GameBranch, "GAME_BRANCH", "public")

	// Apply values with hierarchy
	DiscordToken = getString(cfg.DiscordToken, "DISCORD_TOKEN", "")
	ControlChannelID = getString(cfg.ControlChannelID, "CONTROL_CHANNEL_ID", "")
	StatusChannelID = getString(cfg.StatusChannelID, "STATUS_CHANNEL_ID", "")
	ConnectionListChannelID = getString(cfg.ConnectionListChannelID, "CONNECTION_LIST_CHANNEL_ID", "")
	LogChannelID = getString(cfg.LogChannelID, "LOG_CHANNEL_ID", "")
	SaveChannelID = getString(cfg.SaveChannelID, "SAVE_CHANNEL_ID", "")
	ControlPanelChannelID = getString(cfg.ControlPanelChannelID, "CONTROL_PANEL_CHANNEL_ID", "")
	DiscordCharBufferSize = getInt(cfg.DiscordCharBufferSize, "DISCORD_CHAR_BUFFER_SIZE", 1000)
	BlackListFilePath = getString(cfg.BlackListFilePath, "BLACKLIST_FILE_PATH", "./Blacklist.txt")

	isDiscordEnabledVal := getBool(cfg.IsDiscordEnabled, "IS_DISCORD_ENABLED", false)
	IsDiscordEnabled = isDiscordEnabledVal
	cfg.IsDiscordEnabled = &isDiscordEnabledVal

	ErrorChannelID = getString(cfg.ErrorChannelID, "ERROR_CHANNEL_ID", "")
	BackupKeepLastN = getInt(cfg.BackupKeepLastN, "BACKUP_KEEP_LAST_N", 2000)

	isCleanupEnabledVal := getBool(cfg.IsCleanupEnabled, "IS_CLEANUP_ENABLED", false)
	IsCleanupEnabled = isCleanupEnabledVal
	cfg.IsCleanupEnabled = &isCleanupEnabledVal

	BackupKeepDailyFor = time.Duration(getInt(cfg.BackupKeepDailyFor, "BACKUP_KEEP_DAILY_FOR", 24)) * time.Hour
	BackupKeepWeeklyFor = time.Duration(getInt(cfg.BackupKeepWeeklyFor, "BACKUP_KEEP_WEEKLY_FOR", 168)) * time.Hour
	BackupKeepMonthlyFor = time.Duration(getInt(cfg.BackupKeepMonthlyFor, "BACKUP_KEEP_MONTHLY_FOR", 730)) * time.Hour
	BackupCleanupInterval = time.Duration(getInt(cfg.BackupCleanupInterval, "BACKUP_CLEANUP_INTERVAL", 730)) * time.Hour
	BackupWaitTime = time.Duration(getInt(cfg.BackupWaitTime, "BACKUP_WAIT_TIME", 30)) * time.Second
	Users = getUsers(cfg.Users, "SSUI_USERS", map[string]string{})

	authEnabledVal := getBool(cfg.AuthEnabled, "SSUI_AUTH_ENABLED", false)
	AuthEnabled = authEnabledVal
	cfg.AuthEnabled = &authEnabledVal

	JwtKey = getString(cfg.JwtKey, "SSUI_JWT_KEY", generateJwtKey())
	AuthTokenLifetime = getInt(cfg.AuthTokenLifetime, "SSUI_AUTH_TOKEN_LIFETIME", 1440)

	debugVal := getBool(cfg.Debug, "DEBUG", false)
	IsDebugMode = debugVal
	cfg.Debug = &debugVal

	createSSUILogFileVal := getBool(cfg.CreateSSUILogFile, "CREATE_SSUI_LOGFILE", false)
	CreateSSUILogFile = createSSUILogFileVal
	cfg.CreateSSUILogFile = &createSSUILogFileVal

	LogLevel = getInt(cfg.LogLevel, "LOG_LEVEL", 20)

	isUpdateEnabledVal := getBool(cfg.IsUpdateEnabled, "IS_UPDATE_ENABLED", true)
	IsUpdateEnabled = isUpdateEnabledVal
	cfg.IsUpdateEnabled = &isUpdateEnabledVal

	allowPrereleaseUpdatesVal := getBool(cfg.AllowPrereleaseUpdates, "ALLOW_PRERELEASE_UPDATES", false)
	AllowPrereleaseUpdates = allowPrereleaseUpdatesVal
	cfg.AllowPrereleaseUpdates = &allowPrereleaseUpdatesVal

	allowMajorUpdatesVal := getBool(cfg.AllowMajorUpdates, "ALLOW_MAJOR_UPDATES", false)
	AllowMajorUpdates = allowMajorUpdatesVal
	cfg.AllowMajorUpdates = &allowMajorUpdatesVal

	SubsystemFilters = getStringSlice(cfg.SubsystemFilters, "SUBSYSTEM_FILTERS", []string{})

	isSSCMEnabledVal := getBool(cfg.IsSSCMEnabled, "IS_SSCM_ENABLED", false)
	IsSSCMEnabled = isSSCMEnabledVal
	cfg.IsSSCMEnabled = &isSSCMEnabledVal

	BackendEndpointPort = getString(cfg.BackendEndpointPort, "BACKEND_ENDPOINT_PORT", "8443")
	BackendEndpointIP = getString(cfg.BackendEndpointIP, "BACKEND_ENDPOINT_IP", "0.0.0.0")
	if RunfileGame == "" {
		RunfileGame = getString(cfg.RunfileGame, "RUNFILE_GAME", "Stationeers")
	}
}

// SaveConfig M U S T be called while holding a lock on ConfigMu! Accepts an optional deferred action to run after successfully saving the config
func SaveConfig(deferredAction ...DeferredAction) error {

	cfg := JsonConfig{
		RunfileGame:             RunfileGame,
		BackendEndpointIP:       BackendEndpointIP,
		BackendEndpointPort:     BackendEndpointPort,
		DiscordToken:            DiscordToken,
		ControlChannelID:        ControlChannelID,
		StatusChannelID:         StatusChannelID,
		ConnectionListChannelID: ConnectionListChannelID,
		LogChannelID:            LogChannelID,
		SaveChannelID:           SaveChannelID,
		ControlPanelChannelID:   ControlPanelChannelID,
		DiscordCharBufferSize:   DiscordCharBufferSize,
		BlackListFilePath:       BlackListFilePath,
		IsDiscordEnabled:        &IsDiscordEnabled,
		ErrorChannelID:          ErrorChannelID,
		BackupKeepLastN:         BackupKeepLastN,
		IsCleanupEnabled:        &IsCleanupEnabled,
		BackupKeepDailyFor:      int(BackupKeepDailyFor / time.Hour),
		BackupKeepWeeklyFor:     int(BackupKeepWeeklyFor / time.Hour),
		BackupKeepMonthlyFor:    int(BackupKeepMonthlyFor / time.Hour),
		BackupCleanupInterval:   int(BackupCleanupInterval / time.Hour),
		BackupWaitTime:          int(BackupWaitTime / time.Second),
		GameBranch:              GameBranch,
		Users:                   Users,
		AuthEnabled:             &AuthEnabled,
		JwtKey:                  JwtKey,
		AuthTokenLifetime:       AuthTokenLifetime,
		Debug:                   &IsDebugMode,
		CreateSSUILogFile:       &CreateSSUILogFile,
		LogLevel:                LogLevel,
		SubsystemFilters:        SubsystemFilters,
		IsUpdateEnabled:         &IsUpdateEnabled,
		IsSSCMEnabled:           &IsSSCMEnabled,
		AllowPrereleaseUpdates:  &AllowPrereleaseUpdates,
		AllowMajorUpdates:       &AllowMajorUpdates,
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
