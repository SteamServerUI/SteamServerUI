package backupmgr

import (
	"context"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

func StartBackupLoop() {

	if err := ensureDirectories(); err != nil {
		logger.Backup.Error("Cannot start backup loop: " + err.Error())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if isLoopRunning {
		logger.Backup.Warn("Backup loop is already running")
		return
	}

	ctx, cancel = context.WithCancel(context.Background())
	isLoopRunning = true
	wg.Add(1)

	go backupLoop()
	logger.Backup.Info("Backup loop started")
}

func StopBackupLoop() {
	mu.Lock()
	wasRunning := isLoopRunning
	isLoopRunning = false
	mu.Unlock()

	if !wasRunning {
		return
	}

	cancel()
	wg.Wait()
	logger.Backup.Info("Backup loop stopped")
}

func backupLoop() {
	defer wg.Done()

	ticker := time.NewTicker(cfg.BackupLoopInterval)
	defer ticker.Stop()

	// Perform initial backup
	CreateBackup(cfg.BackupMode)

	for {
		select {
		case <-ctx.Done():
			logger.Backup.Debug("Backup loop cancelled")
			return
		case <-ticker.C:
			CreateBackup(cfg.BackupMode)
		}
	}
}
