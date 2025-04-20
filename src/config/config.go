package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	// All configuration variables can be found in vars.go
	Version = "5.5.17"
	Branch               = "Runfiles"
	IsSteamServerUIBuild = true
)

type JsonConfig struct {
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

type CustomDetection struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Pattern   string `json:"pattern"`
	EventType string `json:"eventType"`
	Message   string `json:"message"`
}

// LoadConfig loads and initializes the configuration
func LoadConfig() (*JsonConfig, error) {
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
		fmt.Println("Config file does not exist. Using defaults and environment variables.")
	} else {
		// Other errors (e.g., permissions), fail immediately
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	// Apply configuration with hierarchy
	applyConfig(&jsonConfig)

	return &jsonConfig, nil
}

// applyConfig applies the configuration with JSON -> env -> fallback hierarchy
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

	// Set backup paths
	//ConfiguredBackupDir = filepath.Join("./saves/", WorldName, "Backup")
	//ConfiguredSafeBackupDir = filepath.Join("./saves/", WorldName, "Safebackups")
}
