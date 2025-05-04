package settings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

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
		{
			Name:        "RunfileGame",
			Type:        "string",
			Group:       "Basic Settings",
			Description: "Runfile Identifier (Restart Required)",
			Value:       config.GetRunfileGame(),
			Required:    true,
		},
		{
			Name:        "IsDebugMode",
			Type:        "bool",
			Group:       "Basic Settings",
			Description: "Enable pprof server",
			Value:       config.GetIsDebugMode(),
		},
		{
			Name:        "CreateSSUILogFile",
			Type:        "bool",
			Group:       "Basic Settings",
			Description: "Create SSUI log files",
			Value:       config.GetCreateSSUILogFile(),
		},
		{
			Name:        "LogLevel",
			Type:        "int",
			Group:       "Basic Settings",
			Description: "Logging verbosity level",
			Value:       config.GetLogLevel(),
			Min:         intPtr(0),
		},
		{
			Name:        "BackendEndpointIP",
			Type:        "string",
			Group:       "Network Settings (Restart Required)",
			Description: "IP address for backend endpoint",
			Value:       config.GetBackendEndpointIP(),
			Required:    true,
		},
		{
			Name:        "BackendEndpointPort",
			Type:        "string",
			Group:       "Network Settings (Restart Required)",
			Description: "Port for backend endpoint",
			Value:       config.GetBackendEndpointPort(),
			Required:    true,
		},
		//{
		//	Name:        "BackupKeepLastN",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Number of recent backups to keep",
		//	Value:       config.GetBackupKeepLastN(),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "IsCleanupEnabled",
		//	Type:        "bool",
		//	Group:       "Advanced Settings",
		//	Description: "Enable backup cleanup",
		//	Value:       config.GetIsCleanupEnabled(),
		//},
		//{
		//	Name:        "BackupKeepDailyFor",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours to keep daily backups",
		//	Value:       int(config.GetBackupKeepDailyFor() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupKeepWeeklyFor",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours to keep weekly backups",
		//	Value:       int(config.GetBackupKeepWeeklyFor() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupKeepMonthlyFor",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours to keep monthly backups",
		//	Value:       int(config.GetBackupKeepMonthlyFor() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupCleanupInterval",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours between backup cleanup runs",
		//	Value:       int(config.GetBackupCleanupInterval() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupWaitTime",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Seconds to wait before backup",
		//	Value:       int(config.GetBackupWaitTime() / time.Second),
		//	Min:         intPtr(0),
		//},
		{
			Name:        "GameBranch",
			Type:        "string",
			Group:       "Advanced Settings",
			Description: "Game branch for updates (Restart Required)",
			Value:       config.GetGameBranch(),
		},
		//{
		//	Name:        "Users",
		//	Type:        "map",
		//	Group:       "Advanced Settings",
		//	Description: "User authentication mappings",
		//	Value:       config.GetUsers(),
		//},
		{
			Name:        "AuthEnabled",
			Type:        "bool",
			Group:       "Advanced Settings",
			Description: "Enable authentication",
			Value:       config.GetAuthEnabled(),
		},
		{
			Name:        "JwtKey",
			Type:        "string",
			Group:       "Advanced Settings",
			Description: "Encryption key for Authentication",
			Value:       config.GetJwtKey(),
		},
		{
			Name:        "AuthTokenLifetime",
			Type:        "int",
			Group:       "Advanced Settings",
			Description: "Token lifetime in seconds",
			Value:       config.GetAuthTokenLifetime(),
			Min:         intPtr(0),
		},
		{
			Name:        "IsUpdateEnabled",
			Type:        "bool",
			Group:       "Advanced Settings",
			Description: "Enable automatic updates",
			Value:       config.GetIsUpdateEnabled(),
		},
		//{
		//	Name:        "IsSSCMEnabled",
		//	Type:        "bool",
		//	Group:       "Advanced Settings",
		//	Description: "Enable SSCM integration",
		//	Value:       config.GetIsSSCMEnabled(),
		//},
		{
			Name:        "AllowPrereleaseUpdates",
			Type:        "bool",
			Group:       "Advanced Settings",
			Description: "Allow prerelease updates",
			Value:       config.GetAllowPrereleaseUpdates(),
		},
		{
			Name:        "AllowMajorUpdates",
			Type:        "bool",
			Group:       "Advanced Settings",
			Description: "Allow major version updates",
			Value:       config.GetAllowMajorUpdates(),
		},
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
