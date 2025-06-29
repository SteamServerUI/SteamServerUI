package backupmgr

import (
	"fmt"
	"os"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

func InitBackupMgr() {
	logger.Backup.Debug("Initializing Backup Manager")

	// Update cfg with current config values
	cfg = Bckupcfg{
		BackupContentDir:   config.GetBackupContentDir(),
		StoredBackupsDir:   config.GetBackupsStoreDir(),
		BackupLoopInterval: config.GetBackupLoopInterval(),
		BackupMode:         config.GetBackupMode(),
		MaxFileSize:        config.GetBackupMaxFileSize(),
		UseCompression:     config.GetBackupUseCompression(),
		KeepSnapshot:       config.GetBackupKeepSnapshot(),
	}

	StartBackupLoop()
	logger.Backup.Info("Content Directory: " + cfg.BackupContentDir)
	logger.Backup.Info("Backup Directory: " + cfg.StoredBackupsDir)
	logger.Backup.Info("Backup Interval: " + cfg.BackupLoopInterval.String())
	logger.Backup.Info("Backup Mode: " + cfg.BackupMode)
	logger.Backup.Info("Max File Size: " + fmt.Sprintf("%d MB", cfg.MaxFileSize/(1024*1024)))
	logger.Backup.Info("Use Compression: " + fmt.Sprintf("%t", cfg.UseCompression))
	logger.Backup.Info("Keep Snapshots: " + fmt.Sprintf("%t", cfg.KeepSnapshot))
}
func ensureDirectories() error {
	dirs := []string{cfg.BackupContentDir, cfg.StoredBackupsDir}
	// if the dirs do not exist, fail and exit without creating them
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("directory %s does not exist", dir)
		}
	}
	return nil
}

func hasContent() bool {
	entries, err := os.ReadDir(cfg.BackupContentDir)
	if err != nil {
		logger.Backup.Warn("Unable to read content directory: " + err.Error())
		return false
	}
	return len(entries) > 0
}
