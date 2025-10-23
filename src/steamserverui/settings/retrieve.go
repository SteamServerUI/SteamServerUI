package settings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

// package settings handles API communication with the config values in package config via getter /setter functions.

// ConfigSetting represents metadata for a configuration setting
type ConfigSetting struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Group       string      `json:"group"`
	Description string      `json:"description"`
	Value       interface{} `json:"value"`
	Min         *int        `json:"min,omitempty"`
	Max         *int        `json:"max,omitempty"`
	Required    bool        `json:"required"`
}

// ConfigSettingsResponse represents the API response
type ConfigSettingsResponse struct {
	Data  []ConfigSetting `json:"data"`
	Error string          `json:"error,omitempty"`
}

func RetrieveSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	settings := []ConfigSetting{
		{Name: "RunfileIdentifier",
			Type:        "string",
			Group:       "Advanced Settings",
			Description: "Runfile identifier. It is recommended to not change this value unsless you know what you are doing.",
			Value:       config.GetRunfileIdentifier(),
		},
		{
			Name:        "IsDebugMode",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Enable pprof server",
			Value:       config.GetIsDebugMode(),
		},
		{
			Name:        "CreateSSUILogFile",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Create SSUI log files",
			Value:       config.GetCreateSSUILogFile(),
		},
		{
			Name:        "LogLevel",
			Type:        "int",
			Group:       "System Settings",
			Description: "Logging verbosity level",
			Value:       config.GetLogLevel(),
			Min:         intPtr(0),
		},
		{
			Name:        "BackendEndpointPort",
			Type:        "string",
			Group:       "System Settings",
			Description: "Port for backend endpoint",
			Value:       config.GetBackendEndpointPort(),
			Required:    true,
		},
		{
			Name:        "GameBranch",
			Type:        "string",
			Group:       "Update Settings",
			Description: "Game branch for updates. Run Steamcmd from the Dashboard after changing this value.",
			Value:       config.GetGameBranch(),
		},
		{
			Name:        "AuthEnabled",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Enable authentication. Be careful, do not lock yourself out!",
			Value:       config.GetAuthEnabled(),
		},
		//{
		//	Name:        "JwtKey",
		//	Type:        "string",
		//	Group:       "Security Settings",
		//	Description: "Encryption key for Authentication",
		//	Value:       config.GetJwtKey(),
		//},
		{
			Name:        "AuthTokenLifetime",
			Type:        "int",
			Group:       "Security Settings",
			Description: "Token lifetime in seconds",
			Value:       config.GetAuthTokenLifetime(),
			Min:         intPtr(0),
		},
		{
			Name:        "IsSSCMEnabled",
			Type:        "bool",
			Group:       "Modding Settings",
			Description: "Enable Unity console hook integration (SteamServerCommandManager) to run commands from the Unity console directly in the game server: Requires BepInEx to be enabled / installed.",
			Value:       config.GetIsSSCMEnabled(),
		},
		{
			Name:        "IsBepInExEnabled",
			Type:        "bool",
			Group:       "Modding Settings",
			Description: "Auto-install BepInEx in the gameserver directory",
			Value:       config.GetIsBepInExEnabled(),
		},
		{
			Name:        "IsSSUICLIConsoleEnabled",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Expose various actions directly in the command line (Restart Required)",
			Value:       config.GetIsSSUICLIConsoleEnabled(),
		},
		{
			Name:        "IsUpdateEnabled",
			Type:        "bool",
			Group:       "Update Settings",
			Description: "Allows automatis SSUI version updates to happen automatically at restart",
			Value:       config.GetIsUpdateEnabled(),
		},
		{
			Name:        "AllowPrereleaseUpdates",
			Type:        "bool",
			Group:       "Update Settings",
			Description: "Allows prerelease SSUI version updates to happen automatically at restart",
			Value:       config.GetAllowPrereleaseUpdates(),
		},
		{
			Name:        "AllowMajorUpdates",
			Type:        "bool",
			Group:       "Update Settings",
			Description: "Allows major SSUI version updates to happen automatically at restart",
			Value:       config.GetAllowMajorUpdates(),
		},
		{
			Name:        "AllowAutoGameServerUpdates",
			Type:        "bool",
			Group:       "Update Settings",
			Description: "Allow automatic game server updates",
			Value:       config.GetAllowAutoGameServerUpdates(),
		},
		{
			Name:        "AutoStartServerOnStartup",
			Type:        "bool",
			Group:       "Advanced Settings",
			Description: "Automatically start the game server on SSUI startup",
			Value:       config.GetAutoStartServerOnStartup(),
		},
		{
			Name:        "AutoRestartServerTimer",
			Type:        "string",
			Group:       "Advanced Settings",
			Description: "Automatically restart the game server after a specified time (e.g., '0' to disable). Suports either timeframe in minutes, or a time like 15:04 or 03:04PM",
			Value:       config.GetAutoRestartServerTimer(),
		},
		{
			Name:        "LogClutterToConsole",
			Type:        "bool",
			Group:       "Advanced Settings",
			Description: "Supresses clutter logs from the gameserver. Useful for Unity servers.",
		},
		// Discord Settings
		{
			Name:        "IsDiscordEnabled",
			Type:        "bool",
			Group:       "Discord Settings",
			Description: "Enable Discord integration",
			Value:       config.GetIsDiscordEnabled(),
		},
		{
			Name:        "DiscordToken",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Discord bot token",
			Value:       config.GetDiscordToken(),
		},
		{
			Name:        "ControlChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Control channel ID",
			Value:       config.GetControlChannelID(),
		},
		{
			Name:        "StatusChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Status channel ID",
			Value:       config.GetStatusChannelID(),
		},
		{
			Name:        "ConnectionListChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Connection list channel ID",
			Value:       config.GetConnectionListChannelID(),
		},
		{
			Name:        "LogChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Log channel ID",
			Value:       config.GetLogChannelID(),
		},
		{
			Name:        "SaveChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Save channel ID",
			Value:       config.GetSaveChannelID(),
		},
		{
			Name:        "ControlPanelChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Control panel channel ID",
			Value:       config.GetControlPanelChannelID(),
		},
		{
			Name:        "DiscordCharBufferSize",
			Type:        "int",
			Group:       "Discord Settings",
			Description: "Discord character buffer size",
			Value:       config.GetDiscordCharBufferSize(),
		},
		{
			Name:        "ErrorChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Error channel ID",
			Value:       config.GetErrorChannelID(),
		},
		//{
		//	Name:        "BackupContentDir",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup content directory",
		//	Value:       config.GetBackupContentDir(),
		//},
		//{
		//	Name:        "BackupsStoreDir",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup stored backups directory",
		//	Value:       config.GetBackupsStoreDir(),
		//},
		//{
		//	Name:        "BackupLoopInterval",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup loop interval",
		//	Value:       config.GetBackupLoopInterval().String(),
		//},
		//{
		//	Name:        "BackupMode",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup mode",
		//	Value:       config.GetBackupMode(),
		//},
		//{
		//	Name:        "BackupMaxFileSize",
		//	Type:        "int",
		//	Group:       "Backup Settings",
		//	Description: "Max file size",
		//	Value:       config.GetBackupMaxFileSize(),
		//},
		//{
		//	Name:        "BackupUseCompression",
		//	Type:        "bool",
		//	Group:       "Backup Settings",
		//	Description: "Use compression",
		//	Value:       config.GetBackupUseCompression(),
		//},
		//{
		//	Name:        "BackupKeepSnapshot",
		//	Type:        "bool",
		//	Group:       "Backup Settings",
		//	Description: "Keep snapshot",
		//	Value:       config.GetBackupKeepSnapshot(),
		//},
	}

	response := ConfigSettingsResponse{
		Data: settings,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// intPtr creates a pointer to an int
func intPtr(i int) *int {
	return &i
}
