package backupmgr

import (
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

// Backup Manager v3 will be a simpler generalized backup system focussing on general backups of a given folder by zipping it as small as possible, then moving it to a backup directory for later restoration

func InitBackupMgr() {
	backupContentDir := "./backupContentDir" //will be changed to config.GetBackupContentDir later which is thread safe, but we leave this for now. This is the directory that will be backed up into zips
	storedBackupsDir := "./storedBackupsDir" //will be changed to config.GetStoredBackupsDir later which is thread safe, but we leave this for now This is the directory that will contain the zip backups
	backupLoopInterval := 5 * time.Minute    //will be changed to config.GetBackupLoopInterval later which is thread safe, but we leave this for now. This is the interval between backups
	backupMode := "zip"                      //will be changed to config.GetBackupMode later which is thread safe, but we leave this for now. This is the mode of backup to use, zip or a plugin name (plugin will follow later)

	//how to use the logger / me trying to make the red lines above not show for this skeleton-commit
	logger.Backup.Debug("Initializing Backup Manager")
	logger.Backup.Info("Backup Manager Initialized")
	logger.Backup.Warn(storedBackupsDir)
	logger.Backup.Info(backupContentDir)
	logger.Backup.Error(backupLoopInterval.String())
	logger.Backup.Error(backupMode)
}
