// backupinterface.go
package backupmgr

import (
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

// GlobalBackupManager is the singleton instance of the backup manager
var GlobalBackupManager *BackupManager

// Track all HTTP handlers that need updating when manager changes
var activeHTTPHandlers []*HTTPHandler

// InitGlobalBackupManager initializes the global backup manager instance
func InitGlobalBackupManager(config BackupConfig) error {
	if GlobalBackupManager != nil {
		GlobalBackupManager.Shutdown()
	}

	GlobalBackupManager = NewBackupManager(config)
	if err := GlobalBackupManager.Initialize(); err != nil {
		return err
	}

	// Update all active HTTP handlers with the new manager
	for _, handler := range activeHTTPHandlers {
		handler.manager = GlobalBackupManager
	}

	return GlobalBackupManager.Start()
}

// RegisterHTTPHandler registers an HTTP handler to be updated when the manager changes
func RegisterHTTPHandler(handler *HTTPHandler) {
	activeHTTPHandlers = append(activeHTTPHandlers, handler)
}

// GetBackupConfig returns a properly configured BackupConfig
func GetBackupConfig() BackupConfig {

	return BackupConfig{
		WorldName:     config.GetWorldName(),
		BackupDir:     config.GetConfiguredBackupDir(),
		SafeBackupDir: config.GetConfiguredSafeBackupDir(),
		WaitTime:      30 * time.Second,
		RetentionPolicy: RetentionPolicy{
			KeepLastN:       config.GetBackupKeepLastN(),
			KeepDailyFor:    config.GetBackupKeepDailyFor(),
			KeepWeeklyFor:   config.GetBackupKeepWeeklyFor(),
			KeepMonthlyFor:  config.GetBackupKeepMonthlyFor(),
			CleanupInterval: config.GetBackupKeepMonthlyFor(),
		},
	}
}

// ReloadBackupManagerFromConfig reloads the global backup manager with the current config. This should be called whenever the config is changed.
func ReloadBackupManagerFromConfig() error {
	// Create a new backupManager config from the global config
	backupConfig := GetBackupConfig()

	// Reinitialize the global backup manager with the new config
	return InitGlobalBackupManager(backupConfig)
}
